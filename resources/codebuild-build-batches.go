package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeBuildBuildBatch struct {
	svc *codebuild.CodeBuild
	Id  *string
}

func init() {
	register("CodeBuildBuildBatch", ListCodeBuildBuildBatch)
}

func ListCodeBuildBuildBatch(sess *session.Session) ([]Resource, error) {
	svc := codebuild.New(sess)
	resources := []Resource{}

	params := &codebuild.ListBuildBatchesInput{}

	for {
		resp, err := svc.ListBuildBatches(params)
		if err != nil {
			return nil, err
		}

		for _, batch := range resp.Ids {
			resources = append(resources, &CodeBuildBuildBatch{
				svc: svc,
				Id:  batch,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeBuildBuildBatch) Remove() error {
	_, err := f.svc.DeleteBuildBatch(&codebuild.DeleteBuildBatchInput{
		Id: f.Id,
	})

	return err
}

func (f *CodeBuildBuildBatch) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Id", f.Id)
	return properties
}
