package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("CloudWatchEventsTarget", ListCloudWatchEventsTargets)
}

func ListCloudWatchEventsTargets(sess *session.Session) ([]Resource, error) {
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
			targetResp, err := svc.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
				Rule:         rule.Name,
				EventBusName: bus.Name,
			})
			if err != nil {
				return nil, err
			}
			for _, target := range targetResp.Targets {
				if err != nil {
					return nil, err
				}

				resources = append(resources, &CloudWatchEventsTarget{
					svc:      svc,
					ruleName: rule.Name,
					targetId: target.Id,
					busName:  bus.Name,
				})
			}
		}
	}

	return resources, nil
}

type CloudWatchEventsTarget struct {
	svc      *cloudwatchevents.CloudWatchEvents
	targetId *string
	ruleName *string
	busName  *string
}

func (target *CloudWatchEventsTarget) Remove() error {
	ids := []*string{target.targetId}
	_, err := target.svc.RemoveTargets(&cloudwatchevents.RemoveTargetsInput{
		Ids:          ids,
		Rule:         target.ruleName,
		EventBusName: target.busName,
		Force:        aws.Bool(true),
	})
	return err
}

func (target *CloudWatchEventsTarget) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("TargetId", target.targetId)
	properties.Set("RuleName", target.ruleName)
	properties.Set("BusName", target.busName)

	return properties
}

func (target *CloudWatchEventsTarget) String() string {
	return fmt.Sprintf("Rule: %s Target ID: %s", *target.ruleName, *target.targetId)
}
