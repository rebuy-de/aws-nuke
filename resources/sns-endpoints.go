package resources

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSEndpoint struct {
	svc *sns.SNS
	ARN *string
}

func init() {
	register("SNSEndpoint", ListSNSEndpoints)
}

func ListSNSEndpoints(sess *session.Session) ([]Resource, error) {
	svc := sns.New(sess)
	resources := []Resource{}
	platformApplications := []*sns.PlatformApplication{}

	platformParams := &sns.ListPlatformApplicationsInput{}

	for {
		resp, err := svc.ListPlatformApplications(platformParams)
		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "InvalidAction" && awsErr.Message() == "Operation (ListPlatformApplications) is not supported in this region" {
				// AWS answers with InvalidAction on regions that do not
				// support ListPlatformApplications.
				break
			}

			return nil, err
		}

		for _, platformApplication := range resp.PlatformApplications {
			platformApplications = append(platformApplications, platformApplication)
		}
		if resp.NextToken == nil {
			break
		}

		platformParams.NextToken = resp.NextToken

	}

	params := &sns.ListEndpointsByPlatformApplicationInput{}

	for _, platformApplication := range platformApplications {

		params.PlatformApplicationArn = platformApplication.PlatformApplicationArn

		resp, err := svc.ListEndpointsByPlatformApplication(params)
		if err != nil {
			return nil, err
		}

		for _, endpoint := range resp.Endpoints {
			resources = append(resources, &SNSEndpoint{
				svc: svc,
				ARN: endpoint.EndpointArn,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *SNSEndpoint) Remove() error {

	_, err := f.svc.DeleteEndpoint(&sns.DeleteEndpointInput{
		EndpointArn: f.ARN,
	})

	return err
}

func (f *SNSEndpoint) String() string {
	return *f.ARN
}
