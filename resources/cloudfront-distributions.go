package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/sirupsen/logrus"
)

type CloudFrontDistribution struct {
	svc    *cloudfront.CloudFront
	ID     *string
	status *string
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
			resources = append(resources, &CloudFrontDistribution{
				svc:    svc,
				ID:     item.Id,
				status: item.Status,
			})
		}

		if *resp.DistributionList.IsTruncated == false {
			break
		}

		params.Marker = resp.DistributionList.NextMarker
	}

	return resources, nil
}

func (f *CloudFrontDistribution) Remove() error {
	logrus.Infof("Waiting on cloudfront id=%s", *f.ID)
	err := f.svc.WaitUntilDistributionDeployed(&cloudfront.GetDistributionInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	// Get Existing eTag
	resp, err := f.svc.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	distributionConfig := resp.DistributionConfig

	if (*distributionConfig.Enabled) {
		distributionConfig.Enabled = new(bool)

		// Cloudfront distribution must be disabled before deleting
		_, err = f.svc.UpdateDistribution(&cloudfront.UpdateDistributionInput{
			Id:                 f.ID,
			IfMatch:            resp.ETag,
			DistributionConfig: distributionConfig,
		})
		if err != nil {
			return err
		}
	}

	logrus.Infof("Waiting on cloudfront id=%s", *f.ID)
	err = f.svc.WaitUntilDistributionDeployed(&cloudfront.GetDistributionInput{
		Id: f.ID,
	})
	if err != nil {
		return err
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