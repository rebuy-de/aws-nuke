package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/clouddirectory"
)

type CloudDirectorySchema struct {
	svc       *clouddirectory.CloudDirectory
	schemaARN *string
}

func init() {
	register("CloudDirectorySchema", ListCloudDirectorySchemas)
}

func ListCloudDirectorySchemas(sess *session.Session) ([]Resource, error) {
	svc := clouddirectory.New(sess)
	resources := []Resource{}

	developmentParams := &clouddirectory.ListDevelopmentSchemaArnsInput{
		MaxResults: aws.Int64(30),
	}

	// Get all development schemas
	for {
		resp, err := svc.ListDevelopmentSchemaArns(developmentParams)
		if err != nil {
			return nil, err
		}

		for _, arn := range resp.SchemaArns {
			resources = append(resources, &CloudDirectorySchema{
				svc:       svc,
				schemaARN: arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		developmentParams.NextToken = resp.NextToken
	}

	// Get all published schemas
	publishedParams := &clouddirectory.ListPublishedSchemaArnsInput{
		MaxResults: aws.Int64(30),
	}
	for {
		resp, err := svc.ListPublishedSchemaArns(publishedParams)
		if err != nil {
			return nil, err
		}

		for _, arn := range resp.SchemaArns {
			resources = append(resources, &CloudDirectorySchema{
				svc:       svc,
				schemaARN: arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		publishedParams.NextToken = resp.NextToken
	}

	// Return combined development and production schemas to DeleteSchema
	return resources, nil
}

func (f *CloudDirectorySchema) Remove() error {

	_, err := f.svc.DeleteSchema(&clouddirectory.DeleteSchemaInput{
		SchemaArn: f.schemaARN,
	})

	return err
}

func (f *CloudDirectorySchema) String() string {
	return *f.schemaARN
}
