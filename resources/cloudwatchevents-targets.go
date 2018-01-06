package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func init() {
	register("CloudWatchEventsTarget", ListCloudWatchEventsTargets)
}

func ListCloudWatchEventsTargets(sess *session.Session) ([]Resource, error) {
	svc := cloudwatchevents.New(sess)

	resp, err := svc.ListRules(nil)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, rule := range resp.Rules {
		targetResp, err := svc.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
			Rule: rule.Name,
		})
		if err != nil {
			return nil, err
		}

		for _, target := range targetResp.Targets {
			resources = append(resources, &CloudWatchEventsTarget{
				svc:      svc,
				ruleName: rule.Name,
				targetId: target.Id,
			})
		}
	}
	return resources, nil
}

type CloudWatchEventsTarget struct {
	svc      *cloudwatchevents.CloudWatchEvents
	targetId *string
	ruleName *string
}

func (target *CloudWatchEventsTarget) Remove() error {
	ids := []*string{target.targetId}
	_, err := target.svc.RemoveTargets(&cloudwatchevents.RemoveTargetsInput{
		Ids:  ids,
		Rule: target.ruleName,
	})
	return err
}

func (target *CloudWatchEventsTarget) String() string {
	return fmt.Sprintf("Rule: %s Target ID: %s", *target.ruleName, *target.targetId)
}
