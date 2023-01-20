package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFv2RuleGroup struct {
	svc       *wafv2.WAFV2
	ID        *string
	name      *string
	lockToken *string
	scope     *string
}

func init() {
	register("WAFv2RuleGroup", ListWAFv2RuleGroups,
		mapCloudControl("AWS::WAFv2::RuleGroup"))
}

func ListWAFv2RuleGroups(sess *session.Session) ([]Resource, error) {
	svc := wafv2.New(sess)
	resources := []Resource{}

	params := &wafv2.ListRuleGroupsInput{
		Limit: aws.Int64(50),
		Scope: aws.String("REGIONAL"),
	}

	output, err := getRuleGroups(svc, params)
	if err != nil {
		return []Resource{}, err
	}

	resources = append(resources, output...)

	if *sess.Config.Region == "us-east-1" {
		params.Scope = aws.String("CLOUDFRONT")

		output, err := getRuleGroups(svc, params)
		if err != nil {
			return []Resource{}, err
		}

		resources = append(resources, output...)
	}

	return resources, nil
}

func getRuleGroups(svc *wafv2.WAFV2, params *wafv2.ListRuleGroupsInput) ([]Resource, error) {
	resources := []Resource{}
	for {
		resp, err := svc.ListRuleGroups(params)
		if err != nil {
			return nil, err
		}

		for _, webACL := range resp.RuleGroups {
			resources = append(resources, &WAFv2RuleGroup{
				svc:       svc,
				ID:        webACL.Id,
				name:      webACL.Name,
				lockToken: webACL.LockToken,
				scope:     params.Scope,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}
	return resources, nil
}

func (f *WAFv2RuleGroup) Remove() error {
	_, err := f.svc.DeleteRuleGroup(&wafv2.DeleteRuleGroupInput{
		Id:        f.ID,
		Name:      f.name,
		Scope:     f.scope,
		LockToken: f.lockToken,
	})

	return err
}

func (f *WAFv2RuleGroup) String() string {
	return *f.ID
}

func (f *WAFv2RuleGroup) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("ID", f.ID).
		Set("Name", f.name).
		Set("Scope", f.scope)
	return properties
}
