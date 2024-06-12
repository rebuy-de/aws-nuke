package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeBuildBuild struct {
	svc *codebuild.CodeBuild
	Id  *string
}

func init() {
	register("CodeBuildBuild", ListCodeBuildBuild)
}

func ListCodeBuildBuild(sess *session.Session) ([]Resource, error) {
	svc := codebuild.New(sess)
	resources := []Resource{}

	params := &codebuild.ListBuildsInput{}

	for {
		resp, err := svc.ListBuilds(params)
		if err != nil {
			return nil, err
		}

		for _, build := range resp.Ids {
			resources = append(resources, &CodeBuildBuild{
				svc: svc,
				Id:  build,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeBuildBuild) Remove() error {
	_, err := f.svc.BatchDeleteBuilds(&codebuild.BatchDeleteBuildsInput{
		Ids: []*string{f.Id},
	})

	return err
}

func (f *CodeBuildBuild) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Id", f.Id)
	return properties
}
