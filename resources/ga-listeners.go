package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

// GlobalAcceleratorListener model
type GlobalAcceleratorListener struct {
	svc *globalaccelerator.GlobalAccelerator
	ARN *string
}

func init() {
	register("GlobalAcceleratorListener", ListGlobalAcceleratorListeners)
}

// ListGlobalAcceleratorListeners enumerates all available listeners of all available accelerators
func ListGlobalAcceleratorListeners(sess *session.Session) ([]Resource, error) {
	svc := globalaccelerator.New(sess)
	acceleratorARNs := []*string{}
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

	// get all listerners
	for _, acceleratorARN := range acceleratorARNs {
		params := &globalaccelerator.ListListenersInput{
			MaxResults:     aws.Int64(100),
			AcceleratorArn: acceleratorARN,
		}

		for {
			output, err := svc.ListListeners(params)
			if err != nil {
				return nil, err
			}

			for _, listener := range output.Listeners {
				resources = append(resources, &GlobalAcceleratorListener{
					svc: svc,
					ARN: listener.ListenerArn,
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
func (gal *GlobalAcceleratorListener) Remove() error {
	_, err := gal.svc.DeleteListener(&globalaccelerator.DeleteListenerInput{
		ListenerArn: gal.ARN,
	})

	return err
}

// Properties definition
func (gal *GlobalAcceleratorListener) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ARN", gal.ARN)
	return properties
}

// String representation
func (gal *GlobalAcceleratorListener) String() string {
	return *gal.ARN
}
