package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontCachePolicy struct {
	svc              *cloudfront.CloudFront
	ID               *string
	Name             *string
	lastModifiedTime *time.Time
}

func init() {
	register("CloudFrontCachePolicy", ListCloudFrontCachePolicies)
}

func ListCloudFrontCachePolicies(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}

	params := &cloudfront.ListCachePoliciesInput{
		MaxItems: aws.Int64(25),
	}

	for {
		resp, err := svc.ListCachePolicies(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.CachePolicyList.Items {
			resources = append(resources, &CloudFrontKeyGroup{
				svc:              svc,
				ID:               item.CachePolicy.Id,
				name:             item.CachePolicy.CachePolicyConfig.Name,
				lastModifiedTime: item.CachePolicy.LastModifiedTime,
			})
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
	properties.Set("Name", f.Name)
	properties.Set("LastModifiedTime", f.lastModifiedTime.Format(time.RFC3339))
	return properties
}
