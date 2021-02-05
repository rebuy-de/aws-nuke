package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type SageMakerNotebookInstanceLifecycleConfig struct {
	svc                 *sagemaker.SageMaker
	lifecycleConfigName *string
}

func init() {
	register("SageMakerNotebookInstanceLifecycleConfig", ListSageMakerNotebookInstanceLifecycleConfigs)
}

func ListSageMakerNotebookInstanceLifecycleConfigs(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListNotebookInstanceLifecycleConfigsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListNotebookInstanceLifecycleConfigs(params)
		if err != nil {
			return nil, err
		}

		for _, lifecycleConfig := range resp.NotebookInstanceLifecycleConfigs {
			resources = append(resources, &SageMakerNotebookInstanceLifecycleConfig{
				svc:                 svc,
				lifecycleConfigName: lifecycleConfig.NotebookInstanceLifecycleConfigName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerNotebookInstanceLifecycleConfig) Remove() error {

	_, err := f.svc.DeleteNotebookInstanceLifecycleConfig(&sagemaker.DeleteNotebookInstanceLifecycleConfigInput{
		NotebookInstanceLifecycleConfigName: f.lifecycleConfigName,
	})

	return err
}

func (f *SageMakerNotebookInstanceLifecycleConfig) String() string {
	return *f.lifecycleConfigName
}

func (f *SageMakerNotebookInstanceLifecycleConfig) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("NotebookInstanceLifecycleConfigName", f.lifecycleConfigName)
	return properties
}
