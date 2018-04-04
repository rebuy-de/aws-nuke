package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTTopicRule struct {
	svc  *iot.IoT
	name *string
}

func init() {
	register("IoTTopicRule", ListIoTTopicRules)
}

func ListIoTTopicRules(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListTopicRulesInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListTopicRules(params)
		if err != nil {
			return nil, err
		}

		for _, rule := range output.Rules {
			resources = append(resources, &IoTTopicRule{
				svc:  svc,
				name: rule.RuleName,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTTopicRule) Remove() error {

	_, err := f.svc.DeleteTopicRule(&iot.DeleteTopicRuleInput{
		RuleName: f.name,
	})

	return err
}

func (f *IoTTopicRule) String() string {
	return *f.name
}
