package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
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

			ruleTagsOutput, err := svc.ListTagsForResource(&cloudwatchevents.ListTagsForResourceInput{
				ResourceARN: rule.Arn,
			})

			if err != nil {
				return nil, err
			}

			resources = append(resources, &CloudWatchEventsRule{
				svc:     svc,
				name:    rule.Name,
				busName: bus.Name,
				tags:    ruleTagsOutput.Tags,
			})
		}
	}
	return resources, nil
}

type CloudWatchEventsRule struct {
	svc     *cloudwatchevents.CloudWatchEvents
	name    *string
	busName *string
	tags    []*cloudwatchevents.Tag
}

func (rule *CloudWatchEventsRule) Remove() error {
	_, err := rule.svc.DeleteRule(&cloudwatchevents.DeleteRuleInput{
		Name:         rule.name,
		EventBusName: rule.busName,
		Force:        aws.Bool(true),
	})
	return err
}

func (rule *CloudWatchEventsRule) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("Name", rule.name)

	for _, tagValue := range rule.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (rule *CloudWatchEventsRule) String() string {
	return fmt.Sprintf("Rule: %s", *rule.name)
}
