package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
)

type CodePipelinePipeline struct {
	svc          *codepipeline.CodePipeline
	pipelineName *string
}

func init() {
	register("CodePipelinePipeline", ListCodePipelinePipelines)
}

func ListCodePipelinePipelines(sess *session.Session) ([]Resource, error) {
	svc := codepipeline.New(sess)
	resources := []Resource{}

	params := &codepipeline.ListPipelinesInput{}

	for {
		resp, err := svc.ListPipelines(params)
		if err != nil {
			return nil, err
		}

		for _, pipeline := range resp.Pipelines {
			resources = append(resources, &CodePipelinePipeline{
				svc:          svc,
				pipelineName: pipeline.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodePipelinePipeline) Remove() error {

	_, err := f.svc.DeletePipeline(&codepipeline.DeletePipelineInput{
		Name: f.pipelineName,
	})

	return err
}

func (f *CodePipelinePipeline) String() string {
	return *f.pipelineName
}
