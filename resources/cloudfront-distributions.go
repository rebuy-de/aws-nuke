package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontDistribution struct {
	svc              *cloudfront.CloudFront
	ID               *string
	status           *string
	lastModifiedTime *time.Time
	tags             []*cloudfront.Tag
}

func init() {
	register("CloudFrontDistribution", ListCloudFrontDistributions)
}

func ListCloudFrontDistributions(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}

	params := &cloudfront.ListDistributionsInput{
		MaxItems: aws.Int64(25),
	}

	for {
		resp, err := svc.ListDistributions(params)
		if err != nil {
			return nil, err
		}
		for _, item := range resp.DistributionList.Items {
			tagResp, err := svc.ListTagsForResource(
				&cloudfront.ListTagsForResourceInput{
					Resource: item.ARN,
				})
			if err != nil {
				return nil, err
			}

			resources = append(resources, &CloudFrontDistribution{
				svc:              svc,
				ID:               item.Id,
				status:           item.Status,
				lastModifiedTime: item.LastModifiedTime,
				tags:             tagResp.Tags.Items,
			})
		}

		if !*resp.DistributionList.IsTruncated {
			break
		}

		params.Marker = resp.DistributionList.NextMarker
	}

	return resources, nil
}

func (f *CloudFrontDistribution) Properties() types.Properties {
	properties := types.NewProperties().
		Set("LastModifiedTime", f.lastModifiedTime.Format(time.RFC3339))

	for _, t := range f.tags {
		properties.SetTag(t.Key, t.Value)
	}
	return properties
}

func (f *CloudFrontDistribution) Remove() error {

	// Get Existing eTag
	resp, err := f.svc.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	if *resp.DistributionConfig.Enabled {
		*resp.DistributionConfig.Enabled = false
		_, err := f.svc.UpdateDistribution(&cloudfront.UpdateDistributionInput{
			Id:                 f.ID,
			DistributionConfig: resp.DistributionConfig,
			IfMatch:            resp.ETag,
		})
		if err != nil {
			return err
		}
	}

	_, err = f.svc.DeleteDistribution(&cloudfront.DeleteDistributionInput{
		Id:      f.ID,
		IfMatch: resp.ETag,
	})

	return err
}

func (f *CloudFrontDistribution) String() string {
	return *f.ID
}
