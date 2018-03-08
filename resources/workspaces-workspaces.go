package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/workspaces"
)

type WorkSpacesWorkspace struct {
	svc         *workspaces.WorkSpaces
	workspaceID *string
}

func init() {
	register("WorkSpacesWorkspace", ListWorkSpacesWorkspaces)
}

func ListWorkSpacesWorkspaces(sess *session.Session) ([]Resource, error) {
	svc := workspaces.New(sess)
	resources := []Resource{}

	params := &workspaces.DescribeWorkspacesInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.DescribeWorkspaces(params)
		if err != nil {
			return nil, err
		}

		for _, workspace := range output.Workspaces {
			resources = append(resources, &WorkSpacesWorkspace{
				svc:         svc,
				workspaceID: workspace.WorkspaceId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *WorkSpacesWorkspace) Remove() error {

	stopRequest := &workspaces.StopRequest{
		WorkspaceId: f.workspaceID,
	}
	terminateRequest := &workspaces.TerminateRequest{
		WorkspaceId: f.workspaceID,
	}
	_, err := f.svc.StopWorkspaces(&workspaces.StopWorkspacesInput{
		StopWorkspaceRequests: []*workspaces.StopRequest{stopRequest},
	})
	if err != nil {
		return err
	}
	_, err = f.svc.TerminateWorkspaces(&workspaces.TerminateWorkspacesInput{
		TerminateWorkspaceRequests: []*workspaces.TerminateRequest{terminateRequest},
	})

	return err
}

func (f *WorkSpacesWorkspace) String() string {
	return *f.workspaceID
}
