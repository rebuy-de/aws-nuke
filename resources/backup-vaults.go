package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type BackupVault struct {
	svc  *backup.Backup
	arn  string
	name string
	tags map[string]*string
}

func init() {
	register("AWSBackupVault", ListBackupVaults)
}

func ListBackupVaults(sess *session.Session) ([]Resource, error) {
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
		tagsOutput, _ := svc.ListTags(&backup.ListTagsInput{ResourceArn: out.BackupVaultArn})
		resources = append(resources, &BackupVault{
			svc:  svc,
			name: *out.BackupVaultName,
			arn:  *out.BackupVaultArn,
			tags: tagsOutput.Tags,
		})
	}

	return resources, nil
}

func (b *BackupVault) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", b.name)
	for tagKey, tagValue := range b.tags {
		properties.Set(fmt.Sprintf("tag:%v", tagKey), *tagValue)
	}
	return properties
}

func (b *BackupVault) Remove() error {
	_, err := b.svc.DeleteBackupVault(&backup.DeleteBackupVaultInput{
		BackupVaultName: &b.name,
	})
	return err
}

func (b *BackupVault) String() string {
	return b.arn
}

func (b *BackupVault) Filter() error {
	if b.name == "aws/efs/automatic-backup-vault" {
		return fmt.Errorf("cannot delete EFS automatic backups vault")
	}
	return nil
}
