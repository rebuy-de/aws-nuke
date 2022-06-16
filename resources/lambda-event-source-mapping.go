package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type LambdaEventSourceMapping struct {
	svc     *lambda.Lambda
	mapping *lambda.EventSourceMappingConfiguration
}

func init() {
	register("LambdaEventSourceMapping", ListLambdaEventSourceMapping)
}

func ListLambdaEventSourceMapping(sess *session.Session) ([]Resource, error) {
	svc := lambda.New(sess)
	resources := []Resource{}

	params := &lambda.ListEventSourceMappingsInput{}
	for {
		resp, err := svc.ListEventSourceMappings(params)
		if err != nil {
			return nil, err
		}

		for _, mapping := range resp.EventSourceMappings {
			resources = append(resources, &LambdaEventSourceMapping{
				svc:     svc,
				mapping: mapping,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.Marker = resp.NextMarker
	}

	return resources, nil
}

func (m *LambdaEventSourceMapping) Remove() error {
	_, err := m.svc.DeleteEventSourceMapping(&lambda.DeleteEventSourceMappingInput{
		UUID: m.mapping.UUID,
	})

	return err
}

func (m *LambdaEventSourceMapping) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("UUID", m.mapping.UUID)
	properties.Set("EventSourceArn", m.mapping.EventSourceArn)
	properties.Set("FunctionArn", m.mapping.FunctionArn)
	properties.Set("State", m.mapping.State)
	return properties
}
