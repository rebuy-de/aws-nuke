package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type BackupRecoveryPoint struct {
	svc             *backup.Backup
	arn             string
	backupVaultName string
}

func init() {
	register("AWSBackupRecoveryPoint", ListBackupRecoveryPoints)
}

func ListBackupRecoveryPoints(sess *session.Session) ([]Resource, error) {
	svc := backup.New(sess)
	max_vaults_len := int64(100)
	params := &backup.ListBackupVaultsInput{
		MaxResults: &max_vaults_len, // aws default limit on number of backup vaults per account
	}
	resp, err := svc.ListBackupVaults(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.BackupVaultList {
		recoveryPointsOutput, _ := svc.ListRecoveryPointsByBackupVault(&backup.ListRecoveryPointsByBackupVaultInput{BackupVaultName: out.BackupVaultName})
		for _, rp := range recoveryPointsOutput.RecoveryPoints {
			resources = append(resources, &BackupRecoveryPoint{
				svc:             svc,
				arn:             *rp.RecoveryPointArn,
				backupVaultName: *out.BackupVaultName,
			})
		}
	}

	return resources, nil
}

func (b *BackupRecoveryPoint) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("BackupVault", b.backupVaultName)
	return properties
}

func (b *BackupRecoveryPoint) Remove() error {
	_, err := b.svc.DeleteRecoveryPoint(&backup.DeleteRecoveryPointInput{
		BackupVaultName:  &b.backupVaultName,
		RecoveryPointArn: &b.arn,
	})
	return err
}

func (b *BackupRecoveryPoint) String() string {
	return fmt.Sprintf("%s", b.arn)
}
