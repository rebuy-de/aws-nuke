package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

type CloudFrontDistributionDeployment struct {
	svc                *cloudfront.CloudFront
	distributionID     *string
	eTag               *string
	distributionConfig *cloudfront.DistributionConfig
	status             string
}

func init() {
	register("CloudFrontDistributionDeployment", ListCloudFrontDistributionDeployments)
}

func ListCloudFrontDistributionDeployments(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}
	distributions := []*cloudfront.DistributionSummary{}

	params := &cloudfront.ListDistributionsInput{
		MaxItems: aws.Int64(25),
	}

	for {
		resp, err := svc.ListDistributions(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.DistributionList.Items {
			distributions = append(distributions, item)
		}

		if *resp.DistributionList.IsTruncated == false {
			break
		}

		params.Marker = resp.DistributionList.NextMarker
	}

	for _, distribution := range distributions {
		params := &cloudfront.GetDistributionInput{
			Id: distribution.Id,
		}
		resp, err := svc.GetDistribution(params)
		if err != nil {
			return nil, err
		}
		resources = append(resources, &CloudFrontDistributionDeployment{
			svc:                svc,
			distributionID:     resp.Distribution.Id,
			eTag:               resp.ETag,
			distributionConfig: resp.Distribution.DistributionConfig,
			status:             UnPtrString(resp.Distribution.Status, "unknown"),
		})
	}

	return resources, nil
}

func (f *CloudFrontDistributionDeployment) Remove() error {

	f.distributionConfig.Enabled = aws.Bool(false)

	_, err := f.svc.UpdateDistribution(&cloudfront.UpdateDistributionInput{
		Id:                 f.distributionID,
		DistributionConfig: f.distributionConfig,
		IfMatch:            f.eTag,
	})

	return err
}

func (f *CloudFrontDistributionDeployment) Filter() error {
	if *f.distributionConfig.Enabled == false && f.status != "InProgress" {
		return fmt.Errorf("already disabled")
	}
	return nil
}

func (f *CloudFrontDistributionDeployment) String() string {
	return *f.distributionID
}
