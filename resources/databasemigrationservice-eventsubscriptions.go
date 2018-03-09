package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

type DatabaseMigrationServiceEventSubscription struct {
	svc              *databasemigrationservice.DatabaseMigrationService
	subscriptionName *string
}

func init() {
	register("DatabaseMigrationServiceEventSubscription", ListDatabaseMigrationServiceEventSubscriptions)
}

func ListDatabaseMigrationServiceEventSubscriptions(sess *session.Session) ([]Resource, error) {
	svc := databasemigrationservice.New(sess)
	resources := []Resource{}

	params := &databasemigrationservice.DescribeEventSubscriptionsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeEventSubscriptions(params)
		if err != nil {
			return nil, err
		}

		for _, eventSubscription := range output.EventSubscriptionsList {
			resources = append(resources, &DatabaseMigrationServiceEventSubscription{
				svc:              svc,
				subscriptionName: eventSubscription.CustSubscriptionId,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *DatabaseMigrationServiceEventSubscription) Remove() error {

	_, err := f.svc.DeleteEventSubscription(&databasemigrationservice.DeleteEventSubscriptionInput{
		SubscriptionName: f.subscriptionName,
	})

	return err
}

func (f *DatabaseMigrationServiceEventSubscription) String() string {
	return *f.subscriptionName
}
