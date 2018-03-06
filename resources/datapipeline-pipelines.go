package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/datapipeline"
)

type DataPipelinePipeline struct {
	svc        *datapipeline.DataPipeline
	pipelineID *string
}

func init() {
	register("DataPipelinePipeline", ListDataPipelinePipelines)
}

func ListDataPipelinePipelines(sess *session.Session) ([]Resource, error) {
	svc := datapipeline.New(sess)
	resources := []Resource{}

	params := &datapipeline.ListPipelinesInput{}

	for {
		resp, err := svc.ListPipelines(params)
		if err != nil {
			return nil, err
		}

		for _, pipeline := range resp.PipelineIdList {
			resources = append(resources, &DataPipelinePipeline{
				svc:        svc,
				pipelineID: pipeline.Id,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (f *DataPipelinePipeline) Remove() error {

	_, err := f.svc.DeletePipeline(&datapipeline.DeletePipelineInput{
		PipelineId: f.pipelineID,
	})

	return err
}

func (f *DataPipelinePipeline) String() string {
	return *f.pipelineID
}
