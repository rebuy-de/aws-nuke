package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworkscm"
)

type OpsWorksCMBackup struct {
	svc *opsworkscm.OpsWorksCM
	ID  *string
}

func init() {
	register("OpsWorksCMBackup", ListOpsWorksCMBackups)
}

func ListOpsWorksCMBackups(sess *session.Session) ([]Resource, error) {
	svc := opsworkscm.New(sess)
	resources := []Resource{}

	params := &opsworkscm.DescribeBackupsInput{}

	output, err := svc.DescribeBackups(params)
	if err != nil {
		return nil, err
	}

	for _, backup := range output.Backups {
		resources = append(resources, &OpsWorksCMBackup{
			svc: svc,
			ID:  backup.BackupId,
		})
	}

	return resources, nil
}

func (f *OpsWorksCMBackup) Remove() error {

	_, err := f.svc.DeleteBackup(&opsworkscm.DeleteBackupInput{
		BackupId: f.ID,
	})

	return err
}

func (f *OpsWorksCMBackup) String() string {
	return *f.ID
}
