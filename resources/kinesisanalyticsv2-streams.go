package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2/kinesisanalyticsv2iface"
)

type KinesisAnalyticsApplicationV2 struct {
	svc             kinesisanalyticsv2iface.KinesisAnalyticsV2API
	applicationName *string
}

func init() {
	register("KinesisAnalyticsApplicationV2", ListKinesisAnalyticsApplicationsV2)
}

func ListKinesisAnalyticsApplicationsV2(sess *session.Session) ([]Resource, error) {
	svc := kinesisanalyticsv2.New(sess)
	resources := []Resource{}
	params := &kinesisanalyticsv2.ListApplicationsInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListApplications(params)
		if err != nil {
			return nil, err
		}

		for _, applicationSummary := range output.ApplicationSummaries {
			resources = append(resources, &KinesisAnalyticsApplicationV2{
				svc:             svc,
				applicationName: applicationSummary.ApplicationName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *KinesisAnalyticsApplicationV2) Remove() error {

	output, err := f.svc.DescribeApplication(&kinesisanalyticsv2.DescribeApplicationInput{
		ApplicationName: f.applicationName,
	})

	if err != nil {
		return err
	}
	createTimestamp := output.ApplicationDetail.CreateTimestamp

	_, err = f.svc.DeleteApplication(&kinesisanalyticsv2.DeleteApplicationInput{
		ApplicationName: f.applicationName,
		CreateTimestamp: createTimestamp,
	})

	return err
}

func (f *KinesisAnalyticsApplicationV2) String() string {
	return *f.applicationName
}
