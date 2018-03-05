package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchLogsDestination struct {
	svc             *cloudwatchlogs.CloudWatchLogs
	destinationName *string
}

func init() {
	register("CloudWatchLogsDestination", ListCloudWatchLogsDestinations)
}

func ListCloudWatchLogsDestinations(sess *session.Session) ([]Resource, error) {
	svc := cloudwatchlogs.New(sess)
	resources := []Resource{}

	params := &cloudwatchlogs.DescribeDestinationsInput{
		Limit: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeDestinations(params)
		if err != nil {
			return nil, err
		}

		for _, destination := range output.Destinations {
			resources = append(resources, &CloudWatchLogsDestination{
				svc:             svc,
				destinationName: destination.DestinationName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CloudWatchLogsDestination) Remove() error {

	_, err := f.svc.DeleteDestination(&cloudwatchlogs.DeleteDestinationInput{
		DestinationName: f.destinationName,
	})

	return err
}

func (f *CloudWatchLogsDestination) String() string {
	return *f.destinationName
}
