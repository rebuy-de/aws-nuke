package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
)

type OpsWorksLayer struct {
	svc *opsworks.OpsWorks
	ID  *string
}

func init() {
	register("OpsWorksLayer", ListOpsWorksLayers)
}

func ListOpsWorksLayers(sess *session.Session) ([]Resource, error) {
	svc := opsworks.New(sess)
	resources := []Resource{}

	stackParams := &opsworks.DescribeStacksInput{}

	resp, err := svc.DescribeStacks(stackParams)
	if err != nil {
		return nil, err
	}

	layerParams := &opsworks.DescribeLayersInput{}

	for _, stack := range resp.Stacks {
		layerParams.StackId = stack.StackId
		output, err := svc.DescribeLayers(layerParams)
		if err != nil {
			return nil, err
		}

		for _, layer := range output.Layers {
			resources = append(resources, &OpsWorksLayer{
				svc: svc,
				ID:  layer.LayerId,
			})
		}
	}

	return resources, nil
}

func (f *OpsWorksLayer) Remove() error {

	_, err := f.svc.DeleteLayer(&opsworks.DeleteLayerInput{
		LayerId: f.ID,
	})

	return err
}

func (f *OpsWorksLayer) String() string {
	return *f.ID
}
