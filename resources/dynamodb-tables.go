package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type DynamoDBTable struct {
	svc  *dynamodb.DynamoDB
	id   string
	tags []*dynamodb.Tag
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
		tags, err := GetTableTags(svc, tableName)

		if err != nil {
			continue
		}

		resources = append(resources, &DynamoDBTable{
			svc: svc,
			id:  *tableName,
			tags: tags,
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

func GetTableTags(svc *dynamodb.DynamoDB, tableName *string) ([]*dynamodb.Tag, error) {
	result, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(*tableName),
	})

	if err != nil {
		return make([]*dynamodb.Tag, 0), err
	}

	tags, err :=  svc.ListTagsOfResource(&dynamodb.ListTagsOfResourceInput{
		ResourceArn: result.Table.TableArn,
	})

	if err != nil {
		return make([]*dynamodb.Tag, 0), err
	}

	return tags.Tags, nil
}

func (i *DynamoDBTable) Properties() types.Properties {
        properties := types.NewProperties()
        properties.Set("Identifier", i.id)

        for _, tag := range i.tags {
                properties.SetTag(tag.Key, tag.Value)
        }

        return properties
}


func (i *DynamoDBTable) String() string {
	return i.id
}
