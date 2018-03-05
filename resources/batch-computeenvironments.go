package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

type BatchComputeEnvironment struct {
	svc                    *batch.Batch
	computeEnvironmentName *string
}

func init() {
	register("BatchComputeEnvironment", ListBatchComputeEnvironments)
}

func ListBatchComputeEnvironments(sess *session.Session) ([]Resource, error) {
	svc := batch.New(sess)
	resources := []Resource{}

	params := &batch.DescribeComputeEnvironmentsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeComputeEnvironments(params)
		if err != nil {
			return nil, err
		}

		for _, computeEnvironment := range output.ComputeEnvironments {
			resources = append(resources, &BatchComputeEnvironment{
				svc: svc,
				computeEnvironmentName: computeEnvironment.ComputeEnvironmentName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *BatchComputeEnvironment) Remove() error {

	_, err := f.svc.DeleteComputeEnvironment(&batch.DeleteComputeEnvironmentInput{
		ComputeEnvironment: f.computeEnvironmentName,
	})

	return err
}

func (f *BatchComputeEnvironment) String() string {
	return *f.computeEnvironmentName
}
