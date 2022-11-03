package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/hunterkepley/aws-nuke/v2/pkg/types"
)

type AMGWorkspace struct {
	svc  *managedgrafana.ManagedGrafana
	id   *string
	name *string
}

func init() {
	register("AMGWorkspace", ListAMGWorkspaces)
}

func ListAMGWorkspaces(sess *session.Session) ([]Resource, error) {
	svc := managedgrafana.New(sess)
	resources := []Resource{}

	var amgWorkspaces []*managedgrafana.WorkspaceSummary
	err := svc.ListWorkspacesPages(
		&managedgrafana.ListWorkspacesInput{},
		func(page *managedgrafana.ListWorkspacesOutput, lastPage bool) bool {
			amgWorkspaces = append(amgWorkspaces, page.Workspaces...)
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	for _, ws := range amgWorkspaces {
		resources = append(resources, &AMGWorkspace{
			svc:  svc,
			id:   ws.Id,
			name: ws.Name,
		})
	}

	return resources, nil
}

func (f *AMGWorkspace) Remove() error {
	_, err := f.svc.DeleteWorkspace(&managedgrafana.DeleteWorkspaceInput{
		WorkspaceId: f.id,
	})

	return err
}

func (f *AMGWorkspace) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("WorkspaceId", f.id).
		Set("WorkspaceName", f.name)

	return properties
}
