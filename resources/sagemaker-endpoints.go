package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

type SageMakerEndpoint struct {
	svc          *sagemaker.SageMaker
	endpointName *string
}

func init() {
	register("SageMakerEndpoint", ListSageMakerEndpoints)
}

func ListSageMakerEndpoints(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListEndpointsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, endpoint := range resp.Endpoints {
			resources = append(resources, &SageMakerEndpoint{
				svc:          svc,
				endpointName: endpoint.EndpointName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerEndpoint) Remove() error {

	_, err := f.svc.DeleteEndpoint(&sagemaker.DeleteEndpointInput{
		EndpointName: f.endpointName,
	})

	return err
}

func (f *SageMakerEndpoint) String() string {
	return *f.endpointName
}
