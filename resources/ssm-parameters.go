package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMParameter struct {
	svc  *ssm.SSM
	name *string
}

func init() {
	register("SSMParameter", ListSSMParameters)
}

func ListSSMParameters(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	params := &ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeParameters(params)
		if err != nil {
			return nil, err
		}

		for _, parameter := range output.Parameters {
			resources = append(resources, &SSMParameter{
				svc:  svc,
				name: parameter.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMParameter) Remove() error {

	_, err := f.svc.DeleteParameter(&ssm.DeleteParameterInput{
		Name: f.name,
	})

	return err
}

func (f *SSMParameter) String() string {
	return *f.name
}
