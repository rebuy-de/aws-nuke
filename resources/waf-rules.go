package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRule struct {
	svc  *waf.WAF
	ID   *string
	rule *waf.Rule
}

func init() {
	register("WAFRule", ListWAFRules)
}

func ListWAFRules(sess *session.Session) ([]Resource, error) {
	svc := waf.New(sess)
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
			ruleResp, _ := svc.GetRule(&waf.GetRuleInput{
				RuleId: rule.RuleId,
			})
			resources = append(resources, &WAFRule{
				svc:  svc,
				ID:   rule.RuleId,
				rule: ruleResp.Rule,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (f *WAFRule) Remove() error {

	tokenOutput, err := f.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	ruleUpdates := []*waf.RuleUpdate{}
	for _, predicate := range f.rule.Predicates {
		ruleUpdates = append(ruleUpdates, &waf.RuleUpdate{
			Action:    aws.String(waf.ChangeActionDelete),
			Predicate: predicate,
		})
	}

	_, err = f.svc.UpdateRule(&waf.UpdateRuleInput{
		ChangeToken: tokenOutput.ChangeToken,
		RuleId:      f.ID,
		Updates:     ruleUpdates,
	})

	if err != nil {
		return err
	}

	tokenOutput, err = f.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteRule(&waf.DeleteRuleInput{
		RuleId:      f.ID,
		ChangeToken: tokenOutput.ChangeToken,
	})

	return err
}

func (f *WAFRule) String() string {
	return *f.ID
}

func (f *WAFRule) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("ID", f.ID).
		Set("Name", f.rule.Name)

	return properties
}
