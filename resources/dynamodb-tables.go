package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBTable struct {
	svc *dynamodb.DynamoDB
	id  string
}

func init() {
	register("DynamoDBTable", ListDynamoDBTables)
}

func ListDynamoDBTables(sess *session.Session) ([]Resource, error) {
	svc := dynamodb.New(sess)

	resp, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, tableName := range resp.TableNames {
		resources = append(resources, &DynamoDBTable{
			svc: svc,
			id:  *tableName,
		})
	}

	return resources, nil
}

func (i *DynamoDBTable) Remove() error {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(i.id),
	}

	_, err := i.svc.DeleteTable(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *DynamoDBTable) String() string {
	return i.id
}
