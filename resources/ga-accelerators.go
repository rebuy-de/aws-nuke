package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

// GAAccelerator model
type GAAccelerator struct {
	svc *globalaccelerator.GlobalAccelerator
	ARN *string
}

func init() {
	register("GAAccelerator", ListGAAccelerators)
}

// ListGAAccelerators enumerates all available accelerators
func ListGAAccelerators(sess *session.Session) ([]Resource, error) {
	svc := globalaccelerator.New(sess)
	resources := []Resource{}

	params := &globalaccelerator.ListAcceleratorsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListAccelerators(params)
		if err != nil {
			return nil, err
		}

		for _, accelerator := range output.Accelerators {
			resources = append(resources, &GAAccelerator{
				svc: svc,
				ARN: accelerator.AcceleratorArn,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

// Remove resource
func (gaa *GAAccelerator) Remove() error {
	_, err := gaa.svc.DeleteAccelerator(&globalaccelerator.DeleteAcceleratorInput{
		AcceleratorArn: gaa.ARN,
	})

	return err
}

// Properties definition
func (gaa *GAAccelerator) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Arn", gaa.ARN)
	return properties
}

// String representation
func (gaa *GAAccelerator) String() string {
	return *gaa.ARN
}
