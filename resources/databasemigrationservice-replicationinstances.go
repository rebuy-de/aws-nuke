package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

type DatabaseMigrationServiceReplicationInstance struct {
	svc *databasemigrationservice.DatabaseMigrationService
	ARN *string
}

func init() {
	register("DatabaseMigrationServiceReplicationInstance", ListDatabaseMigrationServiceReplicationInstances)
}

func ListDatabaseMigrationServiceReplicationInstances(sess *session.Session) ([]Resource, error) {
	svc := databasemigrationservice.New(sess)
	resources := []Resource{}

	params := &databasemigrationservice.DescribeReplicationInstancesInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeReplicationInstances(params)
		if err != nil {
			return nil, err
		}

		for _, replicationInstance := range output.ReplicationInstances {
			resources = append(resources, &DatabaseMigrationServiceReplicationInstance{
				svc: svc,
				ARN: replicationInstance.ReplicationInstanceArn,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *DatabaseMigrationServiceReplicationInstance) Remove() error {

	_, err := f.svc.DeleteReplicationInstance(&databasemigrationservice.DeleteReplicationInstanceInput{
		ReplicationInstanceArn: f.ARN,
	})

	return err
}

func (f *DatabaseMigrationServiceReplicationInstance) String() string {
	return *f.ARN
}
