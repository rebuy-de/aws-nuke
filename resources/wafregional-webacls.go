package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalWebACL struct {
	svc  *wafregional.WAFRegional
	ID   *string
	name *string
}

func init() {
	register("WAFRegionalWebACL", ListWAFRegionalWebACLs)
}

func ListWAFRegionalWebACLs(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListWebACLsInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListWebACLs(params)
		if err != nil {
			return nil, err
		}

		for _, webACL := range resp.WebACLs {
			resources = append(resources, &WAFRegionalWebACL{
				svc:  svc,
				ID:   webACL.WebACLId,
				name: webACL.Name,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (f *WAFRegionalWebACL) Remove() error {

	tokenOutput, err := f.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteWebACL(&waf.DeleteWebACLInput{
		WebACLId:    f.ID,
		ChangeToken: tokenOutput.ChangeToken,
	})

	return err
}

func (f *WAFRegionalWebACL) String() string {
	return *f.ID
}

func (f *WAFRegionalWebACL) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("ID", f.ID).
		Set("Name", f.name)
	return properties
}
