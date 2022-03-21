package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRulePredicate struct {
	svc       *wafregional.WAFRegional
	ruleID    *string
	predicate *waf.Predicate
}

func init() {
	register("WAFRegionalRulePredicate", ListWAFRegionalRulePredicates)
}

func ListWAFRegionalRulePredicates(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListRulesInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListRules(params)
		if err != nil {
			return nil, err
		}

		for _, rule := range resp.Rules {
			details, err := svc.GetRule(&waf.GetRuleInput{
				RuleId: rule.RuleId,
			})
			if err != nil {
				return nil, err
			}

			for _, predicate := range details.Rule.Predicates {
				resources = append(resources, &WAFRegionalRulePredicate{
					svc:       svc,
					ruleID:    rule.RuleId,
					predicate: predicate,
				})
			}
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (r *WAFRegionalRulePredicate) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.UpdateRule(&waf.UpdateRuleInput{
		ChangeToken: tokenOutput.ChangeToken,
		RuleId:      r.ruleID,
		Updates: []*waf.RuleUpdate{
			&waf.RuleUpdate{
				Action:    aws.String("DELETE"),
				Predicate: r.predicate,
			},
		},
	})

	return err
}

func (r *WAFRegionalRulePredicate) Properties() types.Properties {
	return types.NewProperties().
		Set("RuleID", r.ruleID).
		Set("Type", r.predicate.Type).
		Set("Negated", r.predicate.Negated).
		Set("DataID", r.predicate.DataId)
}

func (r *WAFRegionalRulePredicate) String() string {
	return *r.predicate.DataId
}
