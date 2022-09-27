package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudWatchLogsLogGroup struct {
	svc          *cloudwatchlogs.CloudWatchLogs
	logGroupName *string
	tags         map[string]*string
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
			tagParams := &cloudwatchlogs.ListTagsLogGroupInput{
				LogGroupName: logGroup.LogGroupName,
			}

			tagResp, tagErr := svc.ListTagsLogGroup(tagParams)
			if tagErr != nil {
				return nil, tagErr
			}

			resources = append(resources, &CloudWatchLogsLogGroup{
				svc:          svc,
				logGroupName: logGroup.LogGroupName,
				tags:         tagResp.Tags,
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

func (f *CloudWatchLogsLogGroup) String() string {
	return *f.logGroupName
}

func (f *CloudWatchLogsLogGroup) Properties() types.Properties {
	properties := types.NewProperties()
	for k, v := range f.tags {
		properties.SetTag(&k, v)
	}
	properties.
		Set("logGroupName", f.logGroupName)
	return properties
}
