package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

var elbv2ListenerRulePageSize int64 = 400 // AWS has a limit of 100 rules per listener

type ELBv2ListenerRule struct {
	svc         *elbv2.ELBV2
	ruleArn     *string
	lbName      *string
	listenerArn *string
	tags        []*elbv2.Tag
}

func init() {
	register("ELBv2ListenerRule", ListELBv2ListenerRules)
}

func ListELBv2ListenerRules(sess *session.Session) ([]Resource, error) {
	svc := elbv2.New(sess)

	// We need to retrieve ELBs then Listeners then Rules
	lbs := make([]*elbv2.LoadBalancer, 0)
	err := svc.DescribeLoadBalancersPages(
		nil,
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, elbv2 := range page.LoadBalancers {
				lbs = append(lbs, elbv2)
			}
			return !lastPage
		},
	)
	if err != nil {
		return nil, err
	}

	// Required for batched tag retrieval later
	ruleArns := make([]*string, 0)
	ruleArnToResource := make(map[string]*ELBv2ListenerRule)

	resources := make([]Resource, 0)
	for _, lb := range lbs {
		err := svc.DescribeListenersPages(
			&elbv2.DescribeListenersInput{
				LoadBalancerArn: lb.LoadBalancerArn,
			},
			func(page *elbv2.DescribeListenersOutput, lastPage bool) bool {
				for _, listener := range page.Listeners {
					rules, err := svc.DescribeRules(&elbv2.DescribeRulesInput{
						ListenerArn: listener.ListenerArn,
						PageSize:    &elbv2ListenerRulePageSize,
					})
					if err == nil {
						for _, rule := range rules.Rules {
							// Skip default rules as they cannot be deleted
							if rule.IsDefault != nil && *rule.IsDefault {
								continue
							}

							listenerRule := &ELBv2ListenerRule{
								svc:         svc,
								ruleArn:     rule.RuleArn,
								lbName:      lb.LoadBalancerName,
								listenerArn: listener.ListenerArn,
							}

							ruleArns = append(ruleArns, rule.RuleArn)
							resources = append(resources, listenerRule)
							ruleArnToResource[*rule.RuleArn] = listenerRule
						}
					} else {
						logrus.
							WithError(err).
							WithField("listenerArn", listener.ListenerArn).
							Error("Failed to list listener rules for listener")
					}
				}

				return !lastPage
			},
		)
		if err != nil {
			logrus.
				WithError(err).
				WithField("loadBalancerArn", lb.LoadBalancerArn).
				Error("Failed to list listeners for load balancer")
		}
	}

	// Tags for Rules need to be fetched separately
	// We can only specify up to 20 in a single call
	// See: https://github.com/aws/aws-sdk-go/blob/0e8c61841163762f870f6976775800ded4a789b0/service/elbv2/api.go#L5398
	for _, ruleChunk := range Chunk(ruleArns, 20) {
		tagResp, err := svc.DescribeTags(&elbv2.DescribeTagsInput{
			ResourceArns: ruleChunk,
		})
		if err != nil {
			return nil, err
		}
		for _, elbv2TagInfo := range tagResp.TagDescriptions {
			rule := ruleArnToResource[*elbv2TagInfo.ResourceArn]
			rule.tags = elbv2TagInfo.Tags
		}
	}

	return resources, nil
}

func (e *ELBv2ListenerRule) Remove() error {
	_, err := e.svc.DeleteRule(&elbv2.DeleteRuleInput{
		RuleArn: e.ruleArn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *ELBv2ListenerRule) Properties() types.Properties {
	properties := types.NewProperties().
		Set("ARN", e.ruleArn).
		Set("ListenerARN", e.listenerArn).
		Set("LoadBalancerName", e.lbName)

	for _, tagValue := range e.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *ELBv2ListenerRule) String() string {
	return fmt.Sprintf("%s -> %s", *e.lbName, *e.ruleArn)
}
