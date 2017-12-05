package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBTableItem struct {
	svc   *dynamodb.DynamoDB
	id    map[string]*dynamodb.AttributeValue
	table Resource
}

func (n *DynamoDBNuke) ListItems() ([]Resource, error) {
	tables, tablesErr := n.ListTables()
	if tablesErr != nil {
		return nil, tablesErr
	}

	resources := make([]Resource, 0)
	for _, dynamoTable := range tables {
		describeParams := &dynamodb.DescribeTableInput{
			TableName: aws.String(dynamoTable.String()),
		}

		descResp, descErr := n.Service.DescribeTable(describeParams)
		if descErr != nil {
			return nil, descErr
		}

		key := *descResp.Table.KeySchema[0].AttributeName
		params := &dynamodb.ScanInput{
			TableName:            aws.String(dynamoTable.String()),
			ProjectionExpression: aws.String(key),
		}

		scanResp, scanErr := n.Service.Scan(params)
		if scanErr != nil {
			return nil, scanErr
		}

		for _, itemMap := range scanResp.Items {
			resources = append(resources, &DynamoDBTableItem{
				svc:   n.Service,
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
