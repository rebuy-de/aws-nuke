package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRegexMatchTuple struct {
	svc        *wafregional.WAFRegional
	matchSetid *string
	tuple      *waf.RegexMatchTuple
}

func init() {
	register("WAFRegionalRegexMatchTuple", ListWAFRegionalRegexMatchTuple)
}

func ListWAFRegionalRegexMatchTuple(sess *session.Session) ([]Resource, error) {
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
			regexMatchSet, err := svc.GetRegexMatchSet(&waf.GetRegexMatchSetInput{
				RegexMatchSetId: set.RegexMatchSetId,
			})
			if err != nil {
				return nil, err
			}

			for _, tuple := range regexMatchSet.RegexMatchSet.RegexMatchTuples {
				resources = append(resources, &WAFRegionalRegexMatchTuple{
					svc:        svc,
					matchSetid: set.RegexMatchSetId,
					tuple:      tuple,
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

func (r *WAFRegionalRegexMatchTuple) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.UpdateRegexMatchSet(&waf.UpdateRegexMatchSetInput{
		ChangeToken:     tokenOutput.ChangeToken,
		RegexMatchSetId: r.matchSetid,
		Updates: []*waf.RegexMatchSetUpdate{
			&waf.RegexMatchSetUpdate{
				Action:          aws.String("DELETE"),
				RegexMatchTuple: r.tuple,
			},
		},
	})

	return err
}

func (r *WAFRegionalRegexMatchTuple) Properties() types.Properties {
	return types.NewProperties().
		Set("RegexMatchSetID", r.matchSetid).
		Set("FieldToMatchType", r.tuple.FieldToMatch.Type).
		Set("FieldToMatchData", r.tuple.FieldToMatch.Data).
		Set("TextTransformation", r.tuple.TextTransformation)
}
