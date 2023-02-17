package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudWatchAlarm struct {
	svc       *cloudwatch.CloudWatch
	alarmName *string
	tags      []*cloudwatch.Tag
}

func init() {
	register("CloudWatchAlarm", ListCloudWatchAlarms)
}

func ListCloudWatchAlarms(sess *session.Session) ([]Resource, error) {
	svc := cloudwatch.New(sess)
	resources := []Resource{}

	params := &cloudwatch.DescribeAlarmsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeAlarms(params)
		if err != nil {
			return nil, err
		}

		for _, metricAlarm := range output.MetricAlarms {
			tags, err := GetAlarmTags(svc, metricAlarm.AlarmArn)
			if err != nil {
				return nil, err
			}
			resources = append(resources, &CloudWatchAlarm{
				svc:       svc,
				alarmName: metricAlarm.AlarmName,
				tags:      tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func GetAlarmTags(svc *cloudwatch.CloudWatch, arn *string) ([]*cloudwatch.Tag, error) {
	resp, err := svc.ListTagsForResource(&cloudwatch.ListTagsForResourceInput{ResourceARN: arn})
	if err != nil {
		return nil, err
	}

	return resp.Tags, nil
}

func (f *CloudWatchAlarm) Remove() error {

	_, err := f.svc.DeleteAlarms(&cloudwatch.DeleteAlarmsInput{
		AlarmNames: []*string{f.alarmName},
	})

	return err
}

func (f *CloudWatchAlarm) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.alarmName)

	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	return properties
}

func (f *CloudWatchAlarm) String() string {
	return *f.alarmName
}
