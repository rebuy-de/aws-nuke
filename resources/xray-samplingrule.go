package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type XRaySamplingRule struct {
	svc      *xray.XRay
	ruleName *string
	ruleARN  *string
}

func init() {
	register("XRaySamplingRule", ListXRaySamplingRules)
}

func ListXRaySamplingRules(sess *session.Session) ([]Resource, error) {
	svc := xray.New(sess)
	resources := []Resource{}

	var xraySamplingRules []*xray.SamplingRule
	err := svc.GetSamplingRulesPages(
		&xray.GetSamplingRulesInput{},
		func(page *xray.GetSamplingRulesOutput, lastPage bool) bool {
			for _, rule := range page.SamplingRuleRecords {
				if *rule.SamplingRule.RuleName != "Default" {
					xraySamplingRules = append(xraySamplingRules, rule.SamplingRule)
				}
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	for _, rule := range xraySamplingRules {
		resources = append(resources, &XRaySamplingRule{
			svc:      svc,
			ruleName: rule.RuleName,
			ruleARN:  rule.RuleARN,
		})
	}

	return resources, nil
}

func (f *XRaySamplingRule) Remove() error {
	_, err := f.svc.DeleteSamplingRule(&xray.DeleteSamplingRuleInput{
		RuleARN: f.ruleARN, // Specify ruleARN or ruleName, not both
	})

	return err
}

func (f *XRaySamplingRule) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("RuleName", f.ruleName).
		Set("RuleARN", f.ruleARN)

	return properties
}
