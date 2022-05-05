package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalIPSetIP struct {
	svc        *wafregional.WAFRegional
	ipSetid    *string
	descriptor *waf.IPSetDescriptor
}

func init() {
	register("WAFRegionalIPSetIP", ListWAFRegionalIPSetIPs)
}

func ListWAFRegionalIPSetIPs(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListIPSetsInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListIPSets(params)
		if err != nil {
			return nil, err
		}

		for _, set := range resp.IPSets {

			details, err := svc.GetIPSet(&waf.GetIPSetInput{
				IPSetId: set.IPSetId,
			})
			if err != nil {
				return nil, err
			}

			for _, descriptor := range details.IPSet.IPSetDescriptors {
				resources = append(resources, &WAFRegionalIPSetIP{
					svc:        svc,
					ipSetid:    set.IPSetId,
					descriptor: descriptor,
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

func (r *WAFRegionalIPSetIP) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.UpdateIPSet(&waf.UpdateIPSetInput{
		ChangeToken: tokenOutput.ChangeToken,
		IPSetId:     r.ipSetid,
		Updates: []*waf.IPSetUpdate{
			&waf.IPSetUpdate{
				Action:          aws.String("DELETE"),
				IPSetDescriptor: r.descriptor,
			},
		},
	})

	return err
}

func (r *WAFRegionalIPSetIP) Properties() types.Properties {
	return types.NewProperties().
		Set("IPSetID", r.ipSetid).
		Set("Type", r.descriptor.Type).
		Set("Value", r.descriptor.Value)
}
