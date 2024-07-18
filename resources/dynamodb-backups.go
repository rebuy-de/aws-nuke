package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type DynamoDBBackup struct {
	svc *dynamodb.DynamoDB
	id  string
}

func init() {
	register("DynamoDBBackup", ListDynamoDBBackups)
}

func ListDynamoDBBackups(sess *session.Session) ([]Resource, error) {
	svc := dynamodb.New(sess)

	resources := make([]Resource, 0)

	var lastEvaluatedBackupArn *string

	for {
		backupsResp, err := svc.ListBackups(&dynamodb.ListBackupsInput{
			ExclusiveStartBackupArn: lastEvaluatedBackupArn,
		})
		if err != nil {
			return nil, err
		}

		for _, backup := range backupsResp.BackupSummaries {
			resources = append(resources, &DynamoDBBackup{
				svc: svc,
				id:  *backup.BackupArn,
			})
		}

		if backupsResp.LastEvaluatedBackupArn == nil {
			break
		}

		lastEvaluatedBackupArn = backupsResp.LastEvaluatedBackupArn
	}

	return resources, nil
}

func (i *DynamoDBBackup) Remove() error {
	params := &dynamodb.DeleteBackupInput{
		BackupArn: aws.String(i.id),
	}

	_, err := i.svc.DeleteBackup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *DynamoDBBackup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Identifier", i.id)

	return properties
}

func (i *DynamoDBBackup) String() string {
	return i.id
}
