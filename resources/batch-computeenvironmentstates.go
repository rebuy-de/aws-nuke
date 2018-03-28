package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

type BatchComputeEnvironmentState struct {
	svc                    *batch.Batch
	computeEnvironmentName *string
	state                  *string
}

func init() {
	register("BatchComputeEnvironmentState", ListBatchComputeEnvironmentStates)
}

func ListBatchComputeEnvironmentStates(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &BatchComputeEnvironmentState{
				svc: svc,
				computeEnvironmentName: computeEnvironment.ComputeEnvironmentName,
				state: computeEnvironment.State,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *BatchComputeEnvironmentState) Remove() error {

	_, err := f.svc.UpdateComputeEnvironment(&batch.UpdateComputeEnvironmentInput{
		ComputeEnvironment: f.computeEnvironmentName,
		State:              aws.String("DISABLED"),
	})

	return err
}

func (f *BatchComputeEnvironmentState) String() string {
	return *f.computeEnvironmentName
}

func (f *BatchComputeEnvironmentState) Filter() error {
	if strings.ToLower(*f.state) == "disabled" {
		return fmt.Errorf("already disabled")
	}
	return nil
}
