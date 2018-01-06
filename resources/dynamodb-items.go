package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBTableItem struct {
	svc   *dynamodb.DynamoDB
	id    map[string]*dynamodb.AttributeValue
	table Resource
}

func init() {
	register("DynamoDBItem", ListDynamoDBItems)
}

func ListDynamoDBItems(sess *session.Session) ([]Resource, error) {
	svc := dynamodb.New(sess)

	tables, tablesErr := ListDynamoDBTables(sess)
	if tablesErr != nil {
		return nil, tablesErr
	}

	resources := make([]Resource, 0)
	for _, dynamoTable := range tables {
		describeParams := &dynamodb.DescribeTableInput{
			TableName: aws.String(dynamoTable.String()),
		}

		descResp, descErr := svc.DescribeTable(describeParams)
		if descErr != nil {
			return nil, descErr
		}

		key := *descResp.Table.KeySchema[0].AttributeName
		params := &dynamodb.ScanInput{
			TableName:            aws.String(dynamoTable.String()),
			ProjectionExpression: aws.String(key),
		}

		scanResp, scanErr := svc.Scan(params)
		if scanErr != nil {
			return nil, scanErr
		}

		for _, itemMap := range scanResp.Items {
			resources = append(resources, &DynamoDBTableItem{
				svc:   svc,
				id:    itemMap,
				table: dynamoTable,
			})
		}
	}

	return resources, nil
}

func (i *DynamoDBTableItem) Remove() error {
	params := &dynamodb.DeleteItemInput{
		Key:       i.id,
		TableName: aws.String(i.table.String()),
	}

	_, err := i.svc.DeleteItem(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *DynamoDBTableItem) String() string {
	table := i.table.String()
	var keyField string

	for _, value := range i.id {
		value := strings.TrimSpace(value.String())
		keyField = string([]rune(value)[8:(len([]rune(value)) - 3)])
	}

	return table + " -> " + keyField
}
