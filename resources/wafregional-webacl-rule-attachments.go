package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
)

type WAFRegionalWebACLRuleAttachment struct {
	svc           *wafregional.WAFRegional
	webACLID      *string
	activatedRule *waf.ActivatedRule
}

func init() {
	register("WAFRegionalWebACLRuleAttachment", ListWAFRegionalWebACLRuleAttachments)
}

func ListWAFRegionalWebACLRuleAttachments(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}
	webACLs := []*waf.WebACLSummary{}

	params := &waf.ListWebACLsInput{
		Limit: aws.Int64(50),
	}

	//List All Web ACL's
	for {
		resp, err := svc.ListWebACLs(params)
		if err != nil {
			return nil, err
		}

		for _, webACL := range resp.WebACLs {
			webACLs = append(webACLs, webACL)
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	webACLParams := &waf.GetWebACLInput{}

	for _, webACL := range webACLs {
		webACLParams.WebACLId = webACL.WebACLId

		resp, err := svc.GetWebACL(webACLParams)
		if err != nil {
			return nil, err
		}

		for _, webACLRule := range resp.WebACL.Rules {
			resources = append(resources, &WAFRegionalWebACLRuleAttachment{
				svc:           svc,
				webACLID:      webACL.WebACLId,
				activatedRule: webACLRule,
			})
		}

	}

	return resources, nil
}

func (f *WAFRegionalWebACLRuleAttachment) Remove() error {

	tokenOutput, err := f.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	webACLUpdate := &waf.WebACLUpdate{
		Action:        aws.String("DELETE"),
		ActivatedRule: f.activatedRule,
	}

	_, err = f.svc.UpdateWebACL(&waf.UpdateWebACLInput{
		WebACLId:    f.webACLID,
		ChangeToken: tokenOutput.ChangeToken,
		Updates:     []*waf.WebACLUpdate{webACLUpdate},
	})

	return err
}

func (f *WAFRegionalWebACLRuleAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.webACLID, *f.activatedRule.RuleId)
}
