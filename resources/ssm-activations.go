package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMActivation struct {
	svc *ssm.SSM
	ID  *string
}

func init() {
	register("SSMActivation", ListSSMActivations)
}

func ListSSMActivations(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	params := &ssm.DescribeActivationsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeActivations(params)
		if err != nil {
			return nil, err
		}

		for _, activation := range output.ActivationList {
			resources = append(resources, &SSMActivation{
				svc: svc,
				ID:  activation.ActivationId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMActivation) Remove() error {

	_, err := f.svc.DeleteActivation(&ssm.DeleteActivationInput{
		ActivationId: f.ID,
	})

	return err
}

func (f *SSMActivation) String() string {
	return *f.ID
}
