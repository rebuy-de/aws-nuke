package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type DynamoDBTableItem struct {
	svc      *dynamodb.DynamoDB
	id       map[string]*dynamodb.AttributeValue
	table    *DynamoDBTable
	keyName  string
	keyValue string
}

func init() {
	register("DynamoDBTableItem", ListDynamoDBItems)
}

func ListDynamoDBItems(sess *session.Session) ([]Resource, error) {
	svc := dynamodb.New(sess)

	tables, tablesErr := ListDynamoDBTables(sess)
	if tablesErr != nil {
		return nil, tablesErr
	}

	resources := make([]Resource, 0)
	for _, dynamoTableResource := range tables {
		dynamoTable, ok := dynamoTableResource.(*DynamoDBTable)
		if !ok {
			// This should never happen (tm).
			logrus.Errorf("Unable to cast DynamoDBTable.")
			continue
		}

		describeParams := &dynamodb.DescribeTableInput{
			TableName: &dynamoTable.id,
		}

		descResp, descErr := svc.DescribeTable(describeParams)
		if descErr != nil {
			return nil, descErr
		}

		keyName := descResp.Table.KeySchema[0].AttributeName
		params := &dynamodb.ScanInput{
			TableName:            &dynamoTable.id,
			ProjectionExpression: aws.String("#key"),
			ExpressionAttributeNames: map[string]*string{
				"#key": keyName,
			},
		}

		scanResp, scanErr := svc.Scan(params)
		if scanErr != nil {
			return nil, scanErr
		}

		for _, itemMap := range scanResp.Items {
			var keyValue string

			for _, value := range itemMap {
				value := strings.TrimSpace(value.String())
				keyValue = string([]rune(value)[8:(len([]rune(value)) - 3)])
			}

			resources = append(resources, &DynamoDBTableItem{
				svc:      svc,
				id:       itemMap,
				table:    dynamoTable,
				keyName:  aws.StringValue(keyName),
				keyValue: keyValue,
			})
		}
	}

	return resources, nil
}

func (i *DynamoDBTableItem) Remove() error {
	params := &dynamodb.DeleteItemInput{
		Key:       i.id,
		TableName: &i.table.id,
	}

	_, err := i.svc.DeleteItem(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *DynamoDBTableItem) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Table", i.table)
	properties.Set("KeyName", i.keyName)
	properties.Set("KeyValue", i.keyValue)
	return properties
}

func (i *DynamoDBTableItem) String() string {
	return i.table.String() + " -> " + i.keyValue
}
