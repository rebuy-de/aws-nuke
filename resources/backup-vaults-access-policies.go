package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/aws/aws-sdk-go/service/sts"
)

type BackupVaultAccessPolicy struct {
	svc             *backup.Backup
	backupVaultName string
	accountId       string
}

func init() {
	register("AWSBackupVaultAccessPolicy", ListBackupVaultAccessPolicies)
}

func ListBackupVaultAccessPolicies(sess *session.Session) ([]Resource, error) {
	// Lookup current account ID, which is required for building a permissive vault policy
	stsSvc := sts.New(sess)
	identity, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	accountId := *identity.Account

	svc := backup.New(sess)
	maxVaultsLen := int64(100)
	params := &backup.ListBackupVaultsInput{
		MaxResults: &maxVaultsLen, // aws default limit on number of backup vaults per account
	}
	resp, err := svc.ListBackupVaults(params)
	if err != nil {
		return nil, err
	}

	// Iterate over backup vaults and add vault policies that exist.
	resources := make([]Resource, 0)
	for _, out := range resp.BackupVaultList {
		// Check if the Backup Vault has an Access Policy set
		resp, err := svc.GetBackupVaultAccessPolicy(&backup.GetBackupVaultAccessPolicyInput{
			BackupVaultName: out.BackupVaultName,
		})

		// Non-existent Access Policies can come from ResourceNotFoundException or
		// being nil.
		if err != nil {
			switch err.(type) {
			case *backup.ResourceNotFoundException:
				// Non-existent is OK and we skip over them
				continue
			default:
				return nil, err
			}
		}

		// Only delete policies that exist.
		if resp.Policy != nil {
			resources = append(resources, &BackupVaultAccessPolicy{
				svc:             svc,
				backupVaultName: *out.BackupVaultName,
				accountId:       accountId,
			})
		}
	}

	return resources, nil
}

func (b *BackupVaultAccessPolicy) Remove() error {
	// Set the policy to a policy that allows deletion before removal.
	//
	// This is required to delete the policy for the automagically created vaults
	// such as "aws/efs/automatic-backup-vault" from EFS automatic backups
	// which by default Deny policy deletion via backup::DeleteBackupVaultAccessPolicy
	//
	// Example "aws/efs/automatic-backup-vault" default policy:
	//
	// {
	//     "Version": "2012-10-17",
	//     "Statement": [
	//         {
	//             "Effect": "Deny",
	//             "Principal": {
	//                 "AWS": "*"
	//             },
	//             "Action": [
	//                 "backup:DeleteBackupVault",
	//                 "backup:DeleteBackupVaultAccessPolicy",
	//                 "backup:DeleteRecoveryPoint",
	//                 "backup:StartCopyJob",
	//                 "backup:StartRestoreJob",
	//                 "backup:UpdateRecoveryPointLifecycle"
	//             ],
	//             "Resource": "*"
	//         }
	//     ]
	// }
	//
	// While deletion is Denied, you can update the policy with one that
	// doesn't deny and then delete at will.
	allowDeletionPolicy := `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::` + b.accountId + `:root"
            },
            "Action": "backup:DeleteBackupVaultAccessPolicy",
            "Resource": "*"
        }
    ]
}`
	// Ignore error from if we can't put permissive backup vault policy in for some reason, that's OK.
	_, _ = b.svc.PutBackupVaultAccessPolicy(&backup.PutBackupVaultAccessPolicyInput{
		BackupVaultName: &b.backupVaultName,
		Policy:          &allowDeletionPolicy,
	})
	// In the end, this is the only call we actually really care about for err.
	_, err := b.svc.DeleteBackupVaultAccessPolicy(&backup.DeleteBackupVaultAccessPolicyInput{
		BackupVaultName: &b.backupVaultName,
	})
	return err
}

func (b *BackupVaultAccessPolicy) String() string {
	return b.backupVaultName
}
