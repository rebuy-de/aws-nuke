package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

type DatabaseMigrationServiceReplicationTask struct {
	svc *databasemigrationservice.DatabaseMigrationService
	ARN *string
}

func init() {
	register("DatabaseMigrationServiceReplicationTask", ListDatabaseMigrationServiceReplicationTasks)
}

func ListDatabaseMigrationServiceReplicationTasks(sess *session.Session) ([]Resource, error) {
	svc := databasemigrationservice.New(sess)
	resources := []Resource{}

	params := &databasemigrationservice.DescribeReplicationTasksInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeReplicationTasks(params)
		if err != nil {
			return nil, err
		}

		for _, replicationTask := range output.ReplicationTasks {
			resources = append(resources, &DatabaseMigrationServiceReplicationTask{
				svc: svc,
				ARN: replicationTask.ReplicationTaskArn,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *DatabaseMigrationServiceReplicationTask) Remove() error {

	_, err := f.svc.DeleteReplicationTask(&databasemigrationservice.DeleteReplicationTaskInput{
		ReplicationTaskArn: f.ARN,
	})

	return err
}

func (f *DatabaseMigrationServiceReplicationTask) String() string {
	return *f.ARN
}
