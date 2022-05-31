package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type FSxBackup struct {
	svc    *fsx.FSx
	backup *fsx.Backup
}

func init() {
	register("FSxBackup", ListFSxBackups)
}

func ListFSxBackups(sess *session.Session) ([]Resource, error) {
	svc := fsx.New(sess)
	resources := []Resource{}

	params := &fsx.DescribeBackupsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.DescribeBackups(params)
		if err != nil {
			return nil, err
		}

		for _, backup := range resp.Backups {
			resources = append(resources, &FSxBackup{
				svc:    svc,
				backup: backup,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *FSxBackup) Remove() error {
	_, err := f.svc.DeleteBackup(&fsx.DeleteBackupInput{
		BackupId: f.backup.BackupId,
	})

	return err
}

func (f *FSxBackup) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range f.backup.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("Type", f.backup.Type)
	return properties
}

func (f *FSxBackup) String() string {
	return *f.backup.BackupId
}
