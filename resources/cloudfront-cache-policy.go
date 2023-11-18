package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontCachePolicy struct {
	svc *cloudfront.CloudFront
	ID  *string
}

func init() {
	register("CloudFrontCachePolicy", ListCloudFrontCachePolicy)
}

func ListCloudFrontCachePolicy(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}
	params := &cloudfront.ListCachePoliciesInput{}

	for {
		resp, err := svc.ListCachePolicies(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.CachePolicyList.Items {
			if *item.Type == "custom" {
				resources = append(resources, &CloudFrontCachePolicy{
					svc: svc,
					ID:  item.CachePolicy.Id,
				})
			}
		}

		if resp.CachePolicyList.NextMarker == nil {
			break
		}

		params.Marker = resp.CachePolicyList.NextMarker
	}

	return resources, nil
}

func (f *CloudFrontCachePolicy) Remove() error {
	resp, err := f.svc.GetCachePolicy(&cloudfront.GetCachePolicyInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteCachePolicy(&cloudfront.DeleteCachePolicyInput{
		Id:      f.ID,
		IfMatch: resp.ETag,
	})

	return err
}

func (f *CloudFrontCachePolicy) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	return properties
}
