package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRegexPatternString struct {
	svc           *wafregional.WAFRegional
	patternSetid  *string
	patternString *string
}

func init() {
	register("WAFRegionalRegexPatternString", ListWAFRegionalRegexPatternString)
}

func ListWAFRegionalRegexPatternString(sess *session.Session) ([]Resource, error) {
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
			regexPatternSet, err := svc.GetRegexPatternSet(&waf.GetRegexPatternSetInput{
				RegexPatternSetId: set.RegexPatternSetId,
			})
			if err != nil {
				return nil, err
			}

			for _, patternString := range regexPatternSet.RegexPatternSet.RegexPatternStrings {
				resources = append(resources, &WAFRegionalRegexPatternString{
					svc:           svc,
					patternSetid:  set.RegexPatternSetId,
					patternString: patternString,
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

func (r *WAFRegionalRegexPatternString) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.UpdateRegexPatternSet(&waf.UpdateRegexPatternSetInput{
		ChangeToken:       tokenOutput.ChangeToken,
		RegexPatternSetId: r.patternSetid,
		Updates: []*waf.RegexPatternSetUpdate{
			&waf.RegexPatternSetUpdate{
				Action:             aws.String("DELETE"),
				RegexPatternString: r.patternString,
			},
		},
	})

	return err
}

func (r *WAFRegionalRegexPatternString) Properties() types.Properties {
	return types.NewProperties().
		Set("RegexPatternSetID", r.patternSetid).
		Set("patternString", r.patternString)
}
