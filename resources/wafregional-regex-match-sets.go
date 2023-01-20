package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRegexMatchSet struct {
	svc  *wafregional.WAFRegional
	id   *string
	name *string
}

func init() {
	register("WAFRegionalRegexMatchSet", ListWAFRegionalRegexMatchSet)
}

func ListWAFRegionalRegexMatchSet(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListRegexMatchSetsInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListRegexMatchSets(params)
		if err != nil {
			return nil, err
		}

		for _, set := range resp.RegexMatchSets {
			resources = append(resources, &WAFRegionalRegexMatchSet{
				svc:  svc,
				id:   set.RegexMatchSetId,
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

func (r *WAFRegionalRegexMatchSet) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.DeleteRegexMatchSet(&waf.DeleteRegexMatchSetInput{
		RegexMatchSetId: r.id,
		ChangeToken:     tokenOutput.ChangeToken,
	})

	return err
}

func (r *WAFRegionalRegexMatchSet) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", r.id).
		Set("Name", r.name)
}

func (r *WAFRegionalRegexMatchSet) String() string {
	return *r.id
}
