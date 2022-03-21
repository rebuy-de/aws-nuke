package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSEventSubscription struct {
	svc     *rds.RDS
	id      *string
	enabled *bool
	tags    []*rds.Tag
}

func init() {
	register("RDSEventSubscription", ListRDSEventSubscriptions)
}

func ListRDSEventSubscriptions(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeEventSubscriptionsInput{
		MaxRecords: aws.Int64(100),
	}
	resp, err := svc.DescribeEventSubscriptions(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, eventSubscription := range resp.EventSubscriptionsList {
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: eventSubscription.EventSubscriptionArn,
		})

		if err != nil {
			continue
		}

		resources = append(resources, &RDSEventSubscription{
			svc:     svc,
			id:      eventSubscription.CustSubscriptionId,
			enabled: eventSubscription.Enabled,
			tags:    tags.TagList,
		})

	}

	return resources, nil
}

func (i *RDSEventSubscription) Remove() error {
	params := &rds.DeleteEventSubscriptionInput{
		SubscriptionName: i.id,
	}

	_, err := i.svc.DeleteEventSubscription(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSEventSubscription) String() string {
	return *i.id
}

func (i *RDSEventSubscription) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("ID", i.id).
		Set("Enabled", i.enabled)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
