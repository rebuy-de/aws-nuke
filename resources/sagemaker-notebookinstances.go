package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

type SageMakerNotebookInstance struct {
	svc                  *sagemaker.SageMaker
	notebookInstanceName *string
}

func init() {
	register("SageMakerNotebookInstance", ListSageMakerNotebookInstances)
}

func ListSageMakerNotebookInstances(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &SageMakerNotebookInstance{
				svc:                  svc,
				notebookInstanceName: notebookInstance.NotebookInstanceName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerNotebookInstance) Remove() error {
	_, err := f.svc.DeleteNotebookInstance(&sagemaker.DeleteNotebookInstanceInput{
		NotebookInstanceName: f.notebookInstanceName,
	})

	return err
}

func (f *SageMakerNotebookInstance) String() string {
	return *f.notebookInstanceName
}
