package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontOriginRequestPolicy struct {
	svc *cloudfront.CloudFront
	ID  *string
}

func init() {
	register("OriginRequestPolicy", ListCloudFrontOriginRequestPolicies)
}

func ListCloudFrontOriginRequestPolicies(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}
	params := &cloudfront.ListOriginRequestPoliciesInput{}

	for {
		resp, err := svc.ListOriginRequestPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.OriginRequestPolicyList.Items {
			if *item.Type == "custom" {
				resources = append(resources, &CloudFrontOriginRequestPolicy{
					svc: svc,
					ID:  item.OriginRequestPolicy.Id,
				})
			}
		}

		if resp.OriginRequestPolicyList.NextMarker == nil {
			break
		}

		params.Marker = resp.OriginRequestPolicyList.NextMarker
	}

	return resources, nil
}

func (f *CloudFrontOriginRequestPolicy) Remove() error {
	resp, err := f.svc.GetOriginRequestPolicy(&cloudfront.GetOriginRequestPolicyInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteOriginRequestPolicy(&cloudfront.DeleteOriginRequestPolicyInput{
		Id:      f.ID,
		IfMatch: resp.ETag,
	})

	return err
}

func (f *CloudFrontOriginRequestPolicy) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	return properties
}
