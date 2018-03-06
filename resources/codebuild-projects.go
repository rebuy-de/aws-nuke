package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
)

type CodeBuildProject struct {
	svc         *codebuild.CodeBuild
	projectName *string
}

func init() {
	register("CodeBuildProject", ListCodeBuildProjects)
}

func ListCodeBuildProjects(sess *session.Session) ([]Resource, error) {
	svc := codebuild.New(sess)
	resources := []Resource{}

	params := &codebuild.ListProjectsInput{}

	for {
		resp, err := svc.ListProjects(params)
		if err != nil {
			return nil, err
		}

		for _, project := range resp.Projects {
			resources = append(resources, &CodeBuildProject{
				svc:         svc,
				projectName: project,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeBuildProject) Remove() error {

	_, err := f.svc.DeleteProject(&codebuild.DeleteProjectInput{
		Name: f.projectName,
	})

	return err
}

func (f *CodeBuildProject) String() string {
	return *f.projectName
}
