package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

type SageMakerEndpointConfig struct {
	svc                *sagemaker.SageMaker
	endpointConfigName *string
}

func init() {
	register("SageMakerEndpointConfig", ListSageMakerEndpointConfigs)
}

func ListSageMakerEndpointConfigs(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListEndpointConfigsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListEndpointConfigs(params)
		if err != nil {
			return nil, err
		}

		for _, endpointConfig := range resp.EndpointConfigs {
			resources = append(resources, &SageMakerEndpointConfig{
				svc:                svc,
				endpointConfigName: endpointConfig.EndpointConfigName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerEndpointConfig) Remove() error {

	_, err := f.svc.DeleteEndpointConfig(&sagemaker.DeleteEndpointConfigInput{
		EndpointConfigName: f.endpointConfigName,
	})

	return err
}

func (f *SageMakerEndpointConfig) String() string {
	return *f.endpointConfigName
}
