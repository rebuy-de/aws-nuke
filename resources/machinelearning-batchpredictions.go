package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/machinelearning"
	"github.com/sirupsen/logrus"
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
			if aerr, ok := err.(awserr.Error); ok {
				if strings.Contains(aerr.Message(), "AmazonML is no longer available to new customers") {
					logrus.Info("MachineLearningBranchPrediction: AmazonML is no longer available to new customers. Ignore if you haven't set it up.")
					return nil, nil
				}
			}
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
