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

	params := &sns.ListSubscriptionsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListSubscriptions(params)
		if err != nil {
			return nil, err
		}
		for _, subscription := range resp.Subscriptions {
			if *subscription.SubscriptionArn != "PendingConfirmation" {
				resources = append(resources, &SNSSubscription{
					svc:  svc,
					id:   subscription.SubscriptionArn,
					name: subscription.Owner,
				})
			}

		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
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
