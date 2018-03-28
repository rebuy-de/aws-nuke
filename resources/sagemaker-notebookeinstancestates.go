package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

type SageMakerNotebookInstanceState struct {
	svc                  *sagemaker.SageMaker
	notebookInstanceName *string
	instanceStatus       *string
}

func init() {
	register("SageMakerNotebookInstanceState", ListSageMakerNotebookInstanceStates)
}

func ListSageMakerNotebookInstanceStates(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListNotebookInstancesInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListNotebookInstances(params)
		if err != nil {
			return nil, err
		}

		for _, notebookInstance := range resp.NotebookInstances {
			resources = append(resources, &SageMakerNotebookInstanceState{
				svc:                  svc,
				notebookInstanceName: notebookInstance.NotebookInstanceName,
				instanceStatus:       notebookInstance.NotebookInstanceStatus,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerNotebookInstanceState) Remove() error {

	_, err := f.svc.StopNotebookInstance(&sagemaker.StopNotebookInstanceInput{
		NotebookInstanceName: f.notebookInstanceName,
	})

	return err
}

func (f *SageMakerNotebookInstanceState) String() string {
	return *f.notebookInstanceName
}

func (f *SageMakerNotebookInstanceState) Filter() error {
	if strings.ToLower(*f.instanceStatus) == "stopped" {
		return fmt.Errorf("already stopped")
	}
	return nil
}
