package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/machinelearning"
)

type MachineLearningBranchPrediction struct {
	svc *machinelearning.MachineLearning
	ID  *string
}

func init() {
	register("MachineLearningBranchPrediction", ListMachineLearningBranchPredictions)
}

func ListMachineLearningBranchPredictions(sess *session.Session) ([]Resource, error) {
	svc := machinelearning.New(sess)
	resources := []Resource{}

	params := &machinelearning.DescribeBatchPredictionsInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeBatchPredictions(params)
		if err != nil {
			return nil, err
		}

		for _, result := range output.Results {
			resources = append(resources, &MachineLearningBranchPrediction{
				svc: svc,
				ID:  result.BatchPredictionId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MachineLearningBranchPrediction) Remove() error {

	_, err := f.svc.DeleteBatchPrediction(&machinelearning.DeleteBatchPredictionInput{
		BatchPredictionId: f.ID,
	})

	return err
}

func (f *MachineLearningBranchPrediction) String() string {
	return *f.ID
}
