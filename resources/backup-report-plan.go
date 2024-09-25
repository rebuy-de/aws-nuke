package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type BackupReportPlan struct {
	svc            *backup.Backup
	arn            string
	reportPlanName string
}

func init() {
	register("AWSBackupReportPlan", ListBackupReportPlans)
}

func ListBackupReportPlans(sess *session.Session) ([]Resource, error) {
	svc := backup.New(sess)
	max_backups_len := int64(100)
	params := &backup.ListReportPlansInput{
		MaxResults: &max_backups_len, // aws default limit on number of backup plans per account
	}
	resources := make([]Resource, 0)

	for {
		output, err := svc.ListReportPlans(params)
		if err != nil {
			return nil, err
		}

		for _, report := range output.ReportPlans {
			resources = append(resources, &BackupReportPlan{
				svc:            svc,
				arn:            *report.ReportPlanArn,
				reportPlanName: *report.ReportPlanName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (b *BackupReportPlan) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("BackupReport", b.reportPlanName)
	return properties
}

func (b *BackupReportPlan) Remove() error {
	_, err := b.svc.DeleteReportPlan(&backup.DeleteReportPlanInput{
		ReportPlanName: &b.reportPlanName,
	})
	return err
}

func (b *BackupReportPlan) String() string {
	return fmt.Sprintf("%s", b.arn)
}
