package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchLogsResourcePolicy struct {
	svc        *cloudwatchlogs.CloudWatchLogs
	policyName *string
}

func init() {
	register("CloudWatchLogsResourcePolicy", ListCloudWatchLogsResourcePolicies)
}

func ListCloudWatchLogsResourcePolicies(sess *session.Session) ([]Resource, error) {
	svc := cloudwatchlogs.New(sess)
	resources := []Resource{}

	params := &cloudwatchlogs.DescribeResourcePoliciesInput{
		Limit: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeResourcePolicies(params)
		if err != nil {
			return nil, err
		}

		for _, destination := range output.ResourcePolicies {
			resources = append(resources, &CloudWatchLogsResourcePolicy{
				svc:        svc,
				policyName: destination.PolicyName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CloudWatchLogsResourcePolicy) Remove() error {

	_, err := f.svc.DeleteResourcePolicy(&cloudwatchlogs.DeleteResourcePolicyInput{
		PolicyName: f.policyName,
	})

	return err
}

func (f *CloudWatchLogsResourcePolicy) String() string {
	return *f.policyName
}
