package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

type DatabaseMigrationServiceSubnetGroup struct {
	svc *databasemigrationservice.DatabaseMigrationService
	ID  *string
}

func init() {
	register("DatabaseMigrationServiceSubnetGroup", ListDatabaseMigrationServiceSubnetGroups)
}

func ListDatabaseMigrationServiceSubnetGroups(sess *session.Session) ([]Resource, error) {
	svc := databasemigrationservice.New(sess)
	resources := []Resource{}

	params := &databasemigrationservice.DescribeReplicationSubnetGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeReplicationSubnetGroups(params)
		if err != nil {
			return nil, err
		}

		for _, replicationSubnetGroup := range output.ReplicationSubnetGroups {
			resources = append(resources, &DatabaseMigrationServiceSubnetGroup{
				svc: svc,
				ID:  replicationSubnetGroup.ReplicationSubnetGroupIdentifier,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *DatabaseMigrationServiceSubnetGroup) Remove() error {

	_, err := f.svc.DeleteReplicationSubnetGroup(&databasemigrationservice.DeleteReplicationSubnetGroupInput{
		ReplicationSubnetGroupIdentifier: f.ID,
	})

	return err
}

func (f *DatabaseMigrationServiceSubnetGroup) String() string {
	return *f.ID
}
