package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
)

type ElasticTranscoderPipeline struct {
	svc        *elastictranscoder.ElasticTranscoder
	pipelineID *string
}

func init() {
	register("ElasticTranscoderPipeline", ListElasticTranscoderPipelines)
}

func ListElasticTranscoderPipelines(sess *session.Session) ([]Resource, error) {
	svc := elastictranscoder.New(sess)
	resources := []Resource{}

	params := &elastictranscoder.ListPipelinesInput{}

	for {
		resp, err := svc.ListPipelines(params)
		if err != nil {
			return nil, err
		}

		for _, pipeline := range resp.Pipelines {
			resources = append(resources, &ElasticTranscoderPipeline{
				svc:        svc,
				pipelineID: pipeline.Id,
			})
		}

		if resp.NextPageToken == nil {
			break
		}

		params.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ElasticTranscoderPipeline) Remove() error {

	_, err := f.svc.DeletePipeline(&elastictranscoder.DeletePipelineInput{
		Id: f.pipelineID,
	})

	return err
}

func (f *ElasticTranscoderPipeline) String() string {
	return *f.pipelineID
}
