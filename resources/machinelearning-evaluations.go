package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/machinelearning"
)

type MachineLearningEvaluation struct {
	svc *machinelearning.MachineLearning
	ID  *string
}

func init() {
	register("MachineLearningEvaluation", ListMachineLearningEvaluations)
}

func ListMachineLearningEvaluations(sess *session.Session) ([]Resource, error) {
	svc := machinelearning.New(sess)
	resources := []Resource{}

	params := &machinelearning.DescribeEvaluationsInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeEvaluations(params)
		if err != nil {
			return nil, err
		}

		for _, result := range output.Results {
			resources = append(resources, &MachineLearningEvaluation{
				svc: svc,
				ID:  result.EvaluationId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MachineLearningEvaluation) Remove() error {

	_, err := f.svc.DeleteEvaluation(&machinelearning.DeleteEvaluationInput{
		EvaluationId: f.ID,
	})

	return err
}

func (f *MachineLearningEvaluation) String() string {
	return *f.ID
}
