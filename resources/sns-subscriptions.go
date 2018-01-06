package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func init() {
	register("SNSSubscription", ListSNSSubscriptions)
}

func ListSNSSubscriptions(sess *session.Session) ([]Resource, error) {
	svc := sns.New(sess)

	resp, err := svc.ListSubscriptions(nil)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, subscription := range resp.Subscriptions {
		if *subscription.SubscriptionArn != "PendingConfirmation" {
			resources = append(resources, &SNSSubscription{
				svc:  svc,
				id:   subscription.SubscriptionArn,
				name: subscription.Owner,
			})
		}

	}
	return resources, nil
}

type SNSSubscription struct {
	svc  *sns.SNS
	id   *string
	name *string
}

func (subs *SNSSubscription) Remove() error {
	_, err := subs.svc.Unsubscribe(&sns.UnsubscribeInput{
		SubscriptionArn: subs.id,
	})
	return err
}

func (subs *SNSSubscription) String() string {
	return fmt.Sprintf("Owner: %s ARN: %s", *subs.name, *subs.id)
}
