package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

// GlobalAcceleratorEndpointGroup model
type GlobalAcceleratorEndpointGroup struct {
	svc *globalaccelerator.GlobalAccelerator
	ARN *string
}

func init() {
	register("GlobalAcceleratorEndpointGroup", ListGlobalAcceleratorEndpointGroups)
}

// ListGlobalAcceleratorEndpointGroups enumerates all available accelerators
func ListGlobalAcceleratorEndpointGroups(sess *session.Session) ([]Resource, error) {
	svc := globalaccelerator.New(sess)
	acceleratorARNs := []*string{}
	listenerARNs := []*string{}
	resources := []Resource{}

	// get all accelerator arns
	acceleratorParams := &globalaccelerator.ListAcceleratorsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListAccelerators(acceleratorParams)
		if err != nil {
			return nil, err
		}

		for _, accelerator := range output.Accelerators {
			acceleratorARNs = append(acceleratorARNs, accelerator.AcceleratorArn)
		}

		if output.NextToken == nil {
			break
		}

		acceleratorParams.NextToken = output.NextToken
	}

	// get all listerners arns of all accelerators
	for _, acceleratorARN := range acceleratorARNs {
		listenerParams := &globalaccelerator.ListListenersInput{
			MaxResults:     aws.Int64(100),
			AcceleratorArn: acceleratorARN,
		}

		for {
			output, err := svc.ListListeners(listenerParams)
			if err != nil {
				return nil, err
			}

			for _, listener := range output.Listeners {
				listenerARNs = append(listenerARNs, listener.ListenerArn)
			}

			if output.NextToken == nil {
				break
			}

			listenerParams.NextToken = output.NextToken
		}
	}

	// get all endpoints based on all listeners based on all accelerator
	for _, listenerArn := range listenerARNs {
		params := &globalaccelerator.ListEndpointGroupsInput{
			MaxResults:  aws.Int64(100),
			ListenerArn: listenerArn,
		}

		for {
			output, err := svc.ListEndpointGroups(params)
			if err != nil {
				return nil, err
			}

			for _, endpointGroup := range output.EndpointGroups {
				resources = append(resources, &GlobalAcceleratorEndpointGroup{
					svc: svc,
					ARN: endpointGroup.EndpointGroupArn,
				})
			}

			if output.NextToken == nil {
				break
			}

			params.NextToken = output.NextToken
		}
	}

	return resources, nil
}

// Remove resource
func (gaeg *GlobalAcceleratorEndpointGroup) Remove() error {
	_, err := gaeg.svc.DeleteEndpointGroup(&globalaccelerator.DeleteEndpointGroupInput{
		EndpointGroupArn: gaeg.ARN,
	})

	return err
}

// Properties definition
func (gaeg *GlobalAcceleratorEndpointGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ARN", gaeg.ARN)
	return properties
}

// String representation
func (gaeg *GlobalAcceleratorEndpointGroup) String() string {
	return *gaeg.ARN
}
