package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func init() {
	register("CloudWatchEventsRule", ListCloudWatchEventsRules)
}

func ListCloudWatchEventsRules(sess *session.Session) ([]Resource, error) {
	svc := cloudwatchevents.New(sess)

	resp, err := svc.ListEventBuses(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, bus := range resp.EventBuses {
		resp, err := svc.ListRules(&cloudwatchevents.ListRulesInput{
			EventBusName: bus.Name,
		})
		if err != nil {
			return nil, err
		}

		for _, rule := range resp.Rules {
			resources = append(resources, &CloudWatchEventsRule{
				svc:     svc,
				name:    rule.Name,
				busName: bus.Name,
			})
		}
	}
	return resources, nil
}

type CloudWatchEventsRule struct {
	svc     *cloudwatchevents.CloudWatchEvents
	name    *string
	busName *string
}

func (rule *CloudWatchEventsRule) Remove() error {
	_, err := rule.svc.DeleteRule(&cloudwatchevents.DeleteRuleInput{
		Name:         rule.name,
		EventBusName: rule.busName,
		Force:        aws.Bool(true),
	})
	return err
}

func (rule *CloudWatchEventsRule) String() string {
	return fmt.Sprintf("Rule: %s", *rule.name)
}
