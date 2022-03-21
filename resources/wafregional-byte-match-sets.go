package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalByteMatchSet struct {
	svc  *wafregional.WAFRegional
	id   *string
	name *string
}

func init() {
	register("WAFRegionalByteMatchSet", ListWAFRegionalByteMatchSets)
}

func ListWAFRegionalByteMatchSets(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListByteMatchSetsInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListByteMatchSets(params)
		if err != nil {
			return nil, err
		}

		for _, set := range resp.ByteMatchSets {
			resources = append(resources, &WAFRegionalByteMatchSet{
				svc:  svc,
				id:   set.ByteMatchSetId,
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

func (r *WAFRegionalByteMatchSet) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.DeleteByteMatchSet(&waf.DeleteByteMatchSetInput{
		ByteMatchSetId: r.id,
		ChangeToken:    tokenOutput.ChangeToken,
	})

	return err
}

func (r *WAFRegionalByteMatchSet) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", r.id).
		Set("Name", r.name)
}

func (r *WAFRegionalByteMatchSet) String() string {
	return *r.id
}
