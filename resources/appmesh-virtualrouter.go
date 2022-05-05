package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppMeshVirtualRouter struct {
	svc               *appmesh.AppMesh
	meshName          *string
	virtualRouterName *string
}

func init() {
	register("AppMeshVirtualRouter", ListAppMeshVirtualRouters)
}

func ListAppMeshVirtualRouters(sess *session.Session) ([]Resource, error) {
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

	// List VirtualRouters per Mesh
	var vrs []*appmesh.VirtualRouterRef
	for _, meshName := range meshNames {
		err = svc.ListVirtualRoutersPages(
			&appmesh.ListVirtualRoutersInput{
				MeshName: meshName,
			},
			func(page *appmesh.ListVirtualRoutersOutput, lastPage bool) bool {
				for _, vr := range page.VirtualRouters {
					vrs = append(vrs, vr)
				}
				return lastPage
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// Create the resources
	for _, vr := range vrs {
		resources = append(resources, &AppMeshVirtualRouter{
			svc:               svc,
			meshName:          vr.MeshName,
			virtualRouterName: vr.VirtualRouterName,
		})
	}

	return resources, nil
}

func (f *AppMeshVirtualRouter) Remove() error {
	_, err := f.svc.DeleteVirtualRouter(&appmesh.DeleteVirtualRouterInput{
		MeshName:          f.meshName,
		VirtualRouterName: f.virtualRouterName,
	})

	return err
}

func (f *AppMeshVirtualRouter) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("MeshName", f.meshName).
		Set("Name", f.virtualRouterName)

	return properties
}
