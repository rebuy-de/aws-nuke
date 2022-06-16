package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalIPSet struct {
	svc  *wafregional.WAFRegional
	id   *string
	name *string
}

func init() {
	register("WAFRegionalIPSet", ListWAFRegionalIPSets)
}

func ListWAFRegionalIPSets(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &WAFRegionalIPSet{
				svc:  svc,
				id:   set.IPSetId,
				name: set.Name,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (r *WAFRegionalIPSet) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.DeleteIPSet(&waf.DeleteIPSetInput{
		IPSetId:     r.id,
		ChangeToken: tokenOutput.ChangeToken,
	})

	return err
}

func (r *WAFRegionalIPSet) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", r.id).
		Set("Name", r.name)
}

func (r *WAFRegionalIPSet) String() string {
	return *r.id
}
