package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("ComprehendEndpoint", ListComprehendEndpoints)
}

func ListComprehendEndpoints(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListEndpointsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListEndpoints(params)
		if err != nil {
			return nil, err
		}
		for _, endpoint := range resp.EndpointPropertiesList {
			resources = append(resources, &ComprehendEndpoint{
				svc:      svc,
				endpoint: endpoint,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendEndpoint struct {
	svc      *comprehend.Comprehend
	endpoint *comprehend.EndpointProperties
}

func (ce *ComprehendEndpoint) Remove() error {
	_, err := ce.svc.DeleteEndpoint(&comprehend.DeleteEndpointInput{
		EndpointArn: ce.endpoint.EndpointArn,
	})
	return err
}

func (ce *ComprehendEndpoint) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("EndpointArn", ce.endpoint.EndpointArn)
	properties.Set("ModelArn", ce.endpoint.ModelArn)

	return properties
}

func (ce *ComprehendEndpoint) String() string {
	return *ce.endpoint.EndpointArn
}
