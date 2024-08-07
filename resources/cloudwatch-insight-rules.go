package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudWatchInsightRule struct {
	svc  *cloudwatch.CloudWatch
	name *string
}

func init() {
	register("CloudWatchInsightRule", ListCloudWatchInsightRules)
}

func ListCloudWatchInsightRules(sess *session.Session) ([]Resource, error) {
	svc := cloudwatch.New(sess)
	resources := []Resource{}

	params := &cloudwatch.DescribeInsightRulesInput{
		MaxResults: aws.Int64(25),
	}

	for {
		output, err := svc.DescribeInsightRules(params)
		if err != nil {
			return nil, err
		}

		for _, rules := range output.InsightRules {
			resources = append(resources, &CloudWatchInsightRule{
				svc:  svc,
				name: rules.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CloudWatchInsightRule) Remove() error {
	_, err := f.svc.DeleteInsightRules(&cloudwatch.DeleteInsightRulesInput{
		RuleNames: []*string{f.name},
	})

	return err
}

func (f *CloudWatchInsightRule) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.name)
	return properties
}

func (f *CloudWatchInsightRule) String() string {
	return *f.name
}
