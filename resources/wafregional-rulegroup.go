package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRuleGroup struct {
	svc  *wafregional.WAFRegional
	ID   *string
	name *string
}

func init() {
	register("WAFRegionalRuleGroup", ListWAFRegionalRuleGroups)
}

func ListWAFRegionalRuleGroups(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListRuleGroupsInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListRuleGroups(params)
		if err != nil {
			return nil, err
		}

		for _, rule := range resp.RuleGroups {
			resources = append(resources, &WAFRegionalRuleGroup{
				svc:  svc,
				ID:   rule.RuleGroupId,
				name: rule.Name,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (f *WAFRegionalRuleGroup) Remove() error {

	tokenOutput, err := f.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteRuleGroup(&waf.DeleteRuleGroupInput{
		RuleGroupId: f.ID,
		ChangeToken: tokenOutput.ChangeToken,
	})

	return err
}

func (f *WAFRegionalRuleGroup) String() string {
	return *f.ID
}

func (f *WAFRegionalRuleGroup) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("ID", f.ID).
		Set("Name", f.name)
	return properties
}
