package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudWatchAnomalyDetector struct {
	svc      *cloudwatch.CloudWatch
	detector *cloudwatch.AnomalyDetector
}

func init() {
	register("CloudWatchAnomalyDetector", ListCloudWatchAnomalyDetectors)
}

func ListCloudWatchAnomalyDetectors(sess *session.Session) ([]Resource, error) {
	svc := cloudwatch.New(sess)
	resources := []Resource{}

	params := &cloudwatch.DescribeAnomalyDetectorsInput{
		MaxResults: aws.Int64(25),
	}

	for {
		output, err := svc.DescribeAnomalyDetectors(params)
		if err != nil {
			return nil, err
		}

		for _, detector := range output.AnomalyDetectors {
			resources = append(resources, &CloudWatchAnomalyDetector{
				svc:      svc,
				detector: detector,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CloudWatchAnomalyDetector) Remove() error {
	_, err := f.svc.DeleteAnomalyDetector(&cloudwatch.DeleteAnomalyDetectorInput{
		SingleMetricAnomalyDetector: f.detector.SingleMetricAnomalyDetector,
	})

	return err
}

func (f *CloudWatchAnomalyDetector) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("MetricName", f.detector.SingleMetricAnomalyDetector.MetricName)
	return properties
}

func (f *CloudWatchAnomalyDetector) String() string {
	return *f.detector.SingleMetricAnomalyDetector.MetricName
}
