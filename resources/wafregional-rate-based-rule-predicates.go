package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRateBasedRulePredicate struct {
	svc       *wafregional.WAFRegional
	ruleID    *string
	predicate *waf.Predicate
	rateLimit *int64
}

func init() {
	register("WAFRegionalRateBasedRulePredicate", ListWAFRegionalRateBasedRulePredicates)
}

func ListWAFRegionalRateBasedRulePredicates(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListRateBasedRulesInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListRateBasedRules(params)
		if err != nil {
			return nil, err
		}

		for _, rule := range resp.Rules {
			details, err := svc.GetRateBasedRule(&waf.GetRateBasedRuleInput{
				RuleId: rule.RuleId,
			})
			if err != nil {
				return nil, err
			}

			for _, predicate := range details.Rule.MatchPredicates {
				resources = append(resources, &WAFRegionalRateBasedRulePredicate{
					svc:       svc,
					ruleID:    rule.RuleId,
					rateLimit: details.Rule.RateLimit,
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

func (r *WAFRegionalRateBasedRulePredicate) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.UpdateRateBasedRule(&waf.UpdateRateBasedRuleInput{
		ChangeToken: tokenOutput.ChangeToken,
		RuleId:      r.ruleID,
		RateLimit:   r.rateLimit,
		Updates: []*waf.RuleUpdate{
			&waf.RuleUpdate{
				Action:    aws.String("DELETE"),
				Predicate: r.predicate,
			},
		},
	})

	return err
}

func (r *WAFRegionalRateBasedRulePredicate) Properties() types.Properties {
	return types.NewProperties().
		Set("RuleID", r.ruleID).
		Set("Type", r.predicate.Type).
		Set("Negated", r.predicate.Negated).
		Set("DataID", r.predicate.DataId)
}
