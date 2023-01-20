package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppMeshMesh struct {
	svc      *appmesh.AppMesh
	meshName *string
}

func init() {
	register("AppMeshMesh", ListAppMeshMeshes)
}

func ListAppMeshMeshes(sess *session.Session) ([]Resource, error) {
	svc := appmesh.New(sess)
	resources := []Resource{}

	params := &appmesh.ListMeshesInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.ListMeshes(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Meshes {
			resources = append(resources, &AppMeshMesh{
				svc:      svc,
				meshName: item.MeshName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AppMeshMesh) Remove() error {
	_, err := f.svc.DeleteMesh(&appmesh.DeleteMeshInput{
		MeshName: f.meshName,
	})

	return err
}

func (f *AppMeshMesh) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("MeshName", f.meshName)

	return properties
}
