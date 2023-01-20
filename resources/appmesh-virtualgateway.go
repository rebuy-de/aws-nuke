package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppMeshVirtualGateway struct {
	svc                *appmesh.AppMesh
	meshName           *string
	virtualGatewayName *string
}

func init() {
	register("AppMeshVirtualGateway", ListAppMeshVirtualGateways)
}

func ListAppMeshVirtualGateways(sess *session.Session) ([]Resource, error) {
	svc := appmesh.New(sess)
	resources := []Resource{}

	// Get Meshes
	var meshNames []*string
	err := svc.ListMeshesPages(
		&appmesh.ListMeshesInput{},
		func(page *appmesh.ListMeshesOutput, lastPage bool) bool {
			for _, mesh := range page.Meshes {
				meshNames = append(meshNames, mesh.MeshName)
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	// List VirtualGateways per Mesh
	var vgs []*appmesh.VirtualGatewayRef
	for _, meshName := range meshNames {
		err = svc.ListVirtualGatewaysPages(
			&appmesh.ListVirtualGatewaysInput{
				MeshName: meshName,
			},
			func(page *appmesh.ListVirtualGatewaysOutput, lastPage bool) bool {
				for _, vg := range page.VirtualGateways {
					vgs = append(vgs, vg)
				}
				return lastPage
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// Create the resources
	for _, vg := range vgs {
		resources = append(resources, &AppMeshVirtualGateway{
			svc:                svc,
			meshName:           vg.MeshName,
			virtualGatewayName: vg.VirtualGatewayName,
		})
	}

	return resources, nil
}

func (f *AppMeshVirtualGateway) Remove() error {
	_, err := f.svc.DeleteVirtualGateway(&appmesh.DeleteVirtualGatewayInput{
		MeshName:           f.meshName,
		VirtualGatewayName: f.virtualGatewayName,
	})

	return err
}

func (f *AppMeshVirtualGateway) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("MeshName", f.meshName).
		Set("Name", f.virtualGatewayName)

	return properties
}
