package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type BackupPlan struct {
	svc  *backup.Backup
	id   string
	name string
	arn  string
	tags map[string]*string
}

func init() {
	register("AWSBackupPlan", ListBackupPlans)
}

func ListBackupPlans(sess *session.Session) ([]Resource, error) {
	svc := backup.New(sess)
	false_value := false
	max_backups_len := int64(100)
	params := &backup.ListBackupPlansInput{
		IncludeDeleted: &false_value,
		MaxResults:     &max_backups_len, // aws default limit on number of backup plans per account
	}
	resp, err := svc.ListBackupPlans(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.BackupPlansList {
		tagsOutput, _ := svc.ListTags(&backup.ListTagsInput{ResourceArn: out.BackupPlanArn})
		resources = append(resources, &BackupPlan{
			svc:  svc,
			id:   *out.BackupPlanId,
			name: *out.BackupPlanName,
			arn:  *out.BackupPlanArn,
			tags: tagsOutput.Tags,
		})
	}

	return resources, nil
}

func (b *BackupPlan) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", b.id)
	properties.Set("Name", b.name)
	for tagKey, tagValue := range b.tags {
		properties.Set(fmt.Sprintf("tag:%v", tagKey), *tagValue)
	}
	return properties
}

func (b *BackupPlan) Remove() error {
	_, err := b.svc.DeleteBackupPlan(&backup.DeleteBackupPlanInput{
		BackupPlanId: &b.id,
	})
	return err
}

func (b *BackupPlan) String() string {
	return b.arn
}
