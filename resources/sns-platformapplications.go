package resources

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSPlatformApplication struct {
	svc *sns.SNS
	ARN *string
}

func init() {
	register("SNSPlatformApplication", ListSNSPlatformApplications)
}

func ListSNSPlatformApplications(sess *session.Session) ([]Resource, error) {
	svc := sns.New(sess)
	resources := []Resource{}

	params := &sns.ListPlatformApplicationsInput{}

	for {
		resp, err := svc.ListPlatformApplications(params)
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
			resources = append(resources, &SNSPlatformApplication{
				svc: svc,
				ARN: platformApplication.PlatformApplicationArn,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *SNSPlatformApplication) Remove() error {

	_, err := f.svc.DeletePlatformApplication(&sns.DeletePlatformApplicationInput{
		PlatformApplicationArn: f.ARN,
	})

	return err
}

func (f *SNSPlatformApplication) String() string {
	return *f.ARN
}
