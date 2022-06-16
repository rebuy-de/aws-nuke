package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalByteMatchSetIP struct {
	svc        *wafregional.WAFRegional
	matchSetid *string
	tuple      *waf.ByteMatchTuple
}

func init() {
	register("WAFRegionalByteMatchSetIP", ListWAFRegionalByteMatchSetIPs)
}

func ListWAFRegionalByteMatchSetIPs(sess *session.Session) ([]Resource, error) {
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

			details, err := svc.GetByteMatchSet(&waf.GetByteMatchSetInput{
				ByteMatchSetId: set.ByteMatchSetId,
			})
			if err != nil {
				return nil, err
			}

			for _, tuple := range details.ByteMatchSet.ByteMatchTuples {
				resources = append(resources, &WAFRegionalByteMatchSetIP{
					svc:        svc,
					matchSetid: set.ByteMatchSetId,
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

func (r *WAFRegionalByteMatchSetIP) Remove() error {
	tokenOutput, err := r.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = r.svc.UpdateByteMatchSet(&waf.UpdateByteMatchSetInput{
		ChangeToken:    tokenOutput.ChangeToken,
		ByteMatchSetId: r.matchSetid,
		Updates: []*waf.ByteMatchSetUpdate{
			&waf.ByteMatchSetUpdate{
				Action:         aws.String("DELETE"),
				ByteMatchTuple: r.tuple,
			},
		},
	})

	return err
}

func (r *WAFRegionalByteMatchSetIP) Properties() types.Properties {
	return types.NewProperties().
		Set("ByteMatchSetID", r.matchSetid).
		Set("FieldToMatchType", r.tuple.FieldToMatch.Type).
		Set("FieldToMatchData", r.tuple.FieldToMatch.Data).
		Set("TargetString", r.tuple.TargetString)
}
