package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

type DatabaseMigrationServiceEndpoint struct {
	svc *databasemigrationservice.DatabaseMigrationService
	ARN *string
}

func init() {
	register("DatabaseMigrationServiceEndpoint", ListDatabaseMigrationServiceEndpoints)
}

func ListDatabaseMigrationServiceEndpoints(sess *session.Session) ([]Resource, error) {
	svc := databasemigrationservice.New(sess)
	resources := []Resource{}

	params := &databasemigrationservice.DescribeEndpointsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, endpoint := range output.Endpoints {
			resources = append(resources, &DatabaseMigrationServiceEndpoint{
				svc: svc,
				ARN: endpoint.EndpointArn,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *DatabaseMigrationServiceEndpoint) Remove() error {

	_, err := f.svc.DeleteEndpoint(&databasemigrationservice.DeleteEndpointInput{
		EndpointArn: f.ARN,
	})

	return err
}

func (f *DatabaseMigrationServiceEndpoint) String() string {
	return *f.ARN
}
