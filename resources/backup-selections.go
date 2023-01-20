package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"strings"
)

type BackupSelection struct {
	svc           *backup.Backup
	planId        string
	selectionId   string
	selectionName string
}

func init() {
	register("AWSBackupSelection", ListBackupSelections)
}

func ListBackupSelections(sess *session.Session) ([]Resource, error) {
	svc := backup.New(sess)
	false_value := false
	max_backups_len := int64(100)
	params := &backup.ListBackupPlansInput{
		IncludeDeleted: &false_value,
		MaxResults:     &max_backups_len, // aws default limit on number of backup plans per account
	}
	resources := make([]Resource, 0)

	for {
		output, err := svc.ListBackupPlans(params)
		if err != nil {
			return nil, err
		}

		for _, plan := range output.BackupPlansList {
			selectionsOutput, _ := svc.ListBackupSelections(&backup.ListBackupSelectionsInput{BackupPlanId: plan.BackupPlanId})
			for _, selection := range selectionsOutput.BackupSelectionsList {
				resources = append(resources, &BackupSelection{
					svc:           svc,
					planId:        *selection.BackupPlanId,
					selectionId:   *selection.SelectionId,
					selectionName: *selection.SelectionName,
				})
			}
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (b *BackupSelection) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", b.selectionName)
	properties.Set("ID", b.selectionId)
	properties.Set("PlanID", b.planId)
	return properties
}

func (b *BackupSelection) Remove() error {
	_, err := b.svc.DeleteBackupSelection(&backup.DeleteBackupSelectionInput{
		BackupPlanId: &b.planId,
		SelectionId:  &b.selectionId,
	})
	return err
}

func (b *BackupSelection) String() string {
	return fmt.Sprintf("%s (%s)", b.planId, b.selectionId)
}

func (b *BackupSelection) Filter() error {
	if strings.HasPrefix(b.selectionName, "aws/efs/") {
		return fmt.Errorf("cannot delete EFS automatic backups backup selection")
	}
	return nil
}
