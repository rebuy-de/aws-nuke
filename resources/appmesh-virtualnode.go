package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppMeshVirtualNode struct {
	svc             *appmesh.AppMesh
	meshName        *string
	virtualNodeName *string
}

func init() {
	register("AppMeshVirtualNode", ListAppMeshVirtualNodes)
}

func ListAppMeshVirtualNodes(sess *session.Session) ([]Resource, error) {
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

	// List VirtualNodes per Mesh
	var vns []*appmesh.VirtualNodeRef
	for _, meshName := range meshNames {
		err = svc.ListVirtualNodesPages(
			&appmesh.ListVirtualNodesInput{
				MeshName: meshName,
			},
			func(page *appmesh.ListVirtualNodesOutput, lastPage bool) bool {
				for _, vn := range page.VirtualNodes {
					vns = append(vns, vn)
				}
				return lastPage
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// Create the resources
	for _, vn := range vns {
		resources = append(resources, &AppMeshVirtualNode{
			svc:             svc,
			meshName:        vn.MeshName,
			virtualNodeName: vn.VirtualNodeName,
		})
	}

	return resources, nil
}

func (f *AppMeshVirtualNode) Remove() error {
	_, err := f.svc.DeleteVirtualNode(&appmesh.DeleteVirtualNodeInput{
		MeshName:        f.meshName,
		VirtualNodeName: f.virtualNodeName,
	})

	return err
}

func (f *AppMeshVirtualNode) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("MeshName", f.meshName).
		Set("Name", f.virtualNodeName)

	return properties
}
