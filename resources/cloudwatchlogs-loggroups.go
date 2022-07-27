package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchLogsLogGroup struct {
	svc          *cloudwatchlogs.CloudWatchLogs
	logGroupName *string
}

func init() {
	register("CloudWatchLogsLogGroup", ListCloudWatchLogsLogGroups)
}

func ListCloudWatchLogsLogGroups(sess *session.Session) ([]Resource, error) {
	svc := cloudwatchlogs.New(sess)
	resources := []Resource{}

	params := &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeLogGroups(params)
		if err != nil {
			return nil, err
		}

		for _, logGroup := range output.LogGroups {
			resources = append(resources, &CloudWatchLogsLogGroup{
				svc:          svc,
				logGroupName: logGroup.LogGroupName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CloudWatchLogsLogGroup) Remove() error {

	_, err := f.svc.DeleteLogGroup(&cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: f.logGroupName,
	})

	return err
}

func (f *CloudWatchLogsLogGroup) revenant() {}

func (f *CloudWatchLogsLogGroup) String() string {
	return *f.logGroupName
}
