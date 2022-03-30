package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/prometheusservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AMPWorkspace struct {
	svc            *prometheusservice.PrometheusService
	workspaceAlias *string
	workspaceARN   *string
	workspaceId    *string
}

func init() {
	register("AMPWorkspace", ListAMPWorkspaces,
		mapCloudControl("AWS::APS::Workspace"))
}

func ListAMPWorkspaces(sess *session.Session) ([]Resource, error) {
	svc := prometheusservice.New(sess)
	resources := []Resource{}

	var ampWorkspaces []*prometheusservice.WorkspaceSummary
	err := svc.ListWorkspacesPages(
		&prometheusservice.ListWorkspacesInput{},
		func(page *prometheusservice.ListWorkspacesOutput, lastPage bool) bool {
			ampWorkspaces = append(ampWorkspaces, page.Workspaces...)
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	for _, ws := range ampWorkspaces {
		resources = append(resources, &AMPWorkspace{
			svc:            svc,
			workspaceAlias: ws.Alias,
			workspaceARN:   ws.Arn,
			workspaceId:    ws.WorkspaceId,
		})
	}

	return resources, nil
}

func (f *AMPWorkspace) Remove() error {
	_, err := f.svc.DeleteWorkspace(&prometheusservice.DeleteWorkspaceInput{
		WorkspaceId: f.workspaceId,
	})

	return err
}

func (f *AMPWorkspace) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("WorkspaceAlias", f.workspaceAlias).
		Set("WorkspaceARN", f.workspaceARN).
		Set("WorkspaceId", f.workspaceId)

	return properties
}
