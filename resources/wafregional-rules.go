package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
)

type WAFRegionalRule struct {
	svc *wafregional.WAFRegional
	ID  *string
}

func init() {
	register("WAFRegionalRule", ListWAFRegionalRules)
}

func ListWAFRegionalRules(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &WAFRegionalRule{
				svc: svc,
				ID:  rule.RuleId,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (f *WAFRegionalRule) Remove() error {

	tokenOutput, err := f.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteRule(&waf.DeleteRuleInput{
		RuleId:      f.ID,
		ChangeToken: tokenOutput.ChangeToken,
	})

	return err
}

func (f *WAFRegionalRule) String() string {
	return *f.ID
}
