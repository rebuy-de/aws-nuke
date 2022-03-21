package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRegexPatternSet struct {
	svc  *wafregional.WAFRegional
	id   *string
	name *string
}

func init() {
	register("WAFRegionalRegexPatternSet", ListWAFRegionalRegexPatternSet)
}

func ListWAFRegionalRegexPatternSet(sess *session.Session) ([]Resource, error) {
	svc := wafregional.New(sess)
	resources := []Resource{}

	params := &waf.ListRegexPatternSetsInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListRegexPatternSets(params)
		if err != nil {
			return nil, err
		}

		for _, set := range resp.RegexPatternSets {
			resources = append(resources, &WAFRegionalRegexPatternSet{
				svc:  svc,
				id:   set.RegexPatternSetId,
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

func (r *WAFRegionalRegexPatternSet) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.DeleteRegexPatternSet(&waf.DeleteRegexPatternSetInput{
		RegexPatternSetId: r.id,
		ChangeToken:       tokenOutput.ChangeToken,
	})

	return err
}

func (r *WAFRegionalRegexPatternSet) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", r.id).
		Set("Name", r.name)
}

func (r *WAFRegionalRegexPatternSet) String() string {
	return *r.id
}
