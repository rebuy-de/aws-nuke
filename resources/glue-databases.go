package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type GlueDatabase struct {
	svc  *glue.Glue
	name *string
}

func init() {
	register("GlueDatabase", ListGlueDatabases)
}

func ListGlueDatabases(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetDatabasesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.GetDatabases(params)
		if err != nil {
			return nil, err
		}

		for _, database := range output.DatabaseList {
			resources = append(resources, &GlueDatabase{
				svc:  svc,
				name: database.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDatabase) Remove() error {

	_, err := f.svc.DeleteDatabase(&glue.DeleteDatabaseInput{
		Name: f.name,
	})

	return err
}

func (f *GlueDatabase) String() string {
	return *f.name
}
