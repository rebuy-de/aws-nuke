package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
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

func (p *CloudWatchLogsResourcePolicy) Remove() error {

	_, err := p.svc.DeleteResourcePolicy(&cloudwatchlogs.DeleteResourcePolicyInput{
		PolicyName: p.policyName,
	})

	return err
}

func (p *CloudWatchLogsResourcePolicy) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", p.policyName)

	return properties
}

func (p *CloudWatchLogsResourcePolicy) String() string {
	return *p.policyName
}
