package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type GlueDevEndpoint struct {
	svc          *glue.Glue
	endpointName *string
}

func init() {
	register("GlueDevEndpoint", ListGlueDevEndpoints)
}

func ListGlueDevEndpoints(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetDevEndpointsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.GetDevEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, devEndpoint := range output.DevEndpoints {
			resources = append(resources, &GlueDevEndpoint{
				svc:          svc,
				endpointName: devEndpoint.EndpointName,
			})
		}

		// This one API can and does return an empty string
		if output.NextToken == nil || len(*output.NextToken) == 0 {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDevEndpoint) Remove() error {

	_, err := f.svc.DeleteDevEndpoint(&glue.DeleteDevEndpointInput{
		EndpointName: f.endpointName,
	})

	return err
}

func (f *GlueDevEndpoint) String() string {
	return *f.endpointName
}
