package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppMeshVirtualService struct {
	svc                *appmesh.AppMesh
	meshName           *string
	virtualServiceName *string
}

func init() {
	register("AppMeshVirtualService", ListAppMeshVirtualServices)
}

func ListAppMeshVirtualServices(sess *session.Session) ([]Resource, error) {
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

	// List VirtualServices per Mesh
	var vss []*appmesh.VirtualServiceRef
	for _, meshName := range meshNames {
		err = svc.ListVirtualServicesPages(
			&appmesh.ListVirtualServicesInput{
				MeshName: meshName,
			},
			func(page *appmesh.ListVirtualServicesOutput, lastPage bool) bool {
				for _, vs := range page.VirtualServices {
					vss = append(vss, vs)
				}
				return lastPage
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// Create the resources
	for _, vs := range vss {
		resources = append(resources, &AppMeshVirtualService{
			svc:                svc,
			meshName:           vs.MeshName,
			virtualServiceName: vs.VirtualServiceName,
		})
	}

	return resources, nil
}

func (f *AppMeshVirtualService) Remove() error {
	_, err := f.svc.DeleteVirtualService(&appmesh.DeleteVirtualServiceInput{
		MeshName:           f.meshName,
		VirtualServiceName: f.virtualServiceName,
	})

	return err
}

func (f *AppMeshVirtualService) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("MeshName", f.meshName).
		Set("Name", f.virtualServiceName)

	return properties
}
