package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFv2IPSet struct {
	svc       *wafv2.WAFV2
	id        *string
	name      *string
	lockToken *string
	scope     *string
}

func init() {
	register("WAFv2IPSet", ListWAFv2IPSets,
		mapCloudControl("AWS::WAFv2::IPSet"))
}

func ListWAFv2IPSets(sess *session.Session) ([]Resource, error) {
	svc := wafv2.New(sess)
	resources := []Resource{}

	params := &wafv2.ListIPSetsInput{
		Limit: aws.Int64(50),
		Scope: aws.String("REGIONAL"),
	}

	output, err := getIPSets(svc, params)
	if err != nil {
		return []Resource{}, err
	}

	resources = append(resources, output...)

	if *sess.Config.Region == "us-east-1" {
		params.Scope = aws.String("CLOUDFRONT")

		output, err := getIPSets(svc, params)
		if err != nil {
			return []Resource{}, err
		}

		resources = append(resources, output...)
	}

	return resources, nil
}

func getIPSets(svc *wafv2.WAFV2, params *wafv2.ListIPSetsInput) ([]Resource, error) {
	resources := []Resource{}
	for {
		resp, err := svc.ListIPSets(params)
		if err != nil {
			return nil, err
		}

		for _, set := range resp.IPSets {
			resources = append(resources, &WAFv2IPSet{
				svc:       svc,
				id:        set.Id,
				name:      set.Name,
				lockToken: set.LockToken,
				scope:     params.Scope,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (r *WAFv2IPSet) Remove() error {
	_, err := r.svc.DeleteIPSet(&wafv2.DeleteIPSetInput{
		Id:        r.id,
		Name:      r.name,
		Scope:     r.scope,
		LockToken: r.lockToken,
	})

	return err
}

func (r *WAFv2IPSet) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", r.id).
		Set("Name", r.name).
		Set("Scope", r.scope)
}

func (r *WAFv2IPSet) String() string {
	return *r.id
}
