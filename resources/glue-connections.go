package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type GlueConnection struct {
	svc            *glue.Glue
	connectionName *string
}

func init() {
	register("GlueConnection", ListGlueConnections)
}

func ListGlueConnections(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetConnectionsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.GetConnections(params)
		if err != nil {
			return nil, err
		}

		for _, connection := range output.ConnectionList {
			resources = append(resources, &GlueConnection{
				svc:            svc,
				connectionName: connection.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueConnection) Remove() error {

	_, err := f.svc.DeleteConnection(&glue.DeleteConnectionInput{
		ConnectionName: f.connectionName,
	})

	return err
}

func (f *GlueConnection) String() string {
	return *f.connectionName
}
