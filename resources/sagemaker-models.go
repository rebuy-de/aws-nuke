package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

type SageMakerModel struct {
	svc       *sagemaker.SageMaker
	modelName *string
}

func init() {
	register("SageMakerModel", ListSageMakerModels)
}

func ListSageMakerModels(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListModelsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListModels(params)
		if err != nil {
			return nil, err
		}

		for _, model := range resp.Models {
			resources = append(resources, &SageMakerModel{
				svc:       svc,
				modelName: model.ModelName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerModel) Remove() error {

	_, err := f.svc.DeleteModel(&sagemaker.DeleteModelInput{
		ModelName: f.modelName,
	})

	return err
}

func (f *SageMakerModel) String() string {
	return *f.modelName
}
