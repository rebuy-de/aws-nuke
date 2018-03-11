package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesisanalytics"
)

type KinesisAnalyticsApplication struct {
	svc             *kinesisanalytics.KinesisAnalytics
	applicationName *string
}

func init() {
	register("KinesisAnalyticsApplication", ListKinesisAnalyticsApplications)
}

func ListKinesisAnalyticsApplications(sess *session.Session) ([]Resource, error) {
	svc := kinesisanalytics.New(sess)
	resources := []Resource{}
	var lastApplicationName *string
	params := &kinesisanalytics.ListApplicationsInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListApplications(params)
		if err != nil {
			return nil, err
		}

		for _, applicationSummary := range output.ApplicationSummaries {
			resources = append(resources, &KinesisAnalyticsApplication{
				svc:             svc,
				applicationName: applicationSummary.ApplicationName,
			})
			lastApplicationName = applicationSummary.ApplicationName
		}

		if *output.HasMoreApplications == false {
			break
		}

		params.ExclusiveStartApplicationName = lastApplicationName
	}

	return resources, nil
}

func (f *KinesisAnalyticsApplication) Remove() error {

	output, err := f.svc.DescribeApplication(&kinesisanalytics.DescribeApplicationInput{
		ApplicationName: f.applicationName,
	})

	if err != nil {
		return err
	}
	createTimestamp := output.ApplicationDetail.CreateTimestamp

	_, err = f.svc.DeleteApplication(&kinesisanalytics.DeleteApplicationInput{
		ApplicationName: f.applicationName,
		CreateTimestamp: createTimestamp,
	})

	return err
}

func (f *KinesisAnalyticsApplication) String() string {
	return *f.applicationName
}
