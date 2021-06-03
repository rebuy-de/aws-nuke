package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type CodeBuildProject struct {
	svc         *codebuild.CodeBuild
	projectName *string
	tags        map[string]*string
}

func init() {
	register("CodeBuildProject", ListCodeBuildProjects)
}

func GetTags(svc *codebuild.CodeBuild, projectName []*string) map[string]*string {
	tags := make(map[string]*string)
	batchResult, _ := svc.BatchGetProjects(&codebuild.BatchGetProjectsInput{Names: projectName})
	tags[*batchResult.Projects[0].Tags[0].Key] = batchResult.Projects[0].Tags[0].Value
	return tags
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
				tags:        GetTags(svc, resp.Projects),
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

func (f *CodeBuildProject) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range f.tags {
		properties.SetTag(&key, tag)
	}
	properties.
		Set("Project Name", f.projectName)
	return properties
}
