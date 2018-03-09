package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/machinelearning"
)

type MachineLearningDataSource struct {
	svc *machinelearning.MachineLearning
	ID  *string
}

func init() {
	register("MachineLearningDataSource", ListMachineLearningDataSources)
}

func ListMachineLearningDataSources(sess *session.Session) ([]Resource, error) {
	svc := machinelearning.New(sess)
	resources := []Resource{}

	params := &machinelearning.DescribeDataSourcesInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeDataSources(params)
		if err != nil {
			return nil, err
		}

		for _, result := range output.Results {
			resources = append(resources, &MachineLearningDataSource{
				svc: svc,
				ID:  result.DataSourceId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MachineLearningDataSource) Remove() error {

	_, err := f.svc.DeleteDataSource(&machinelearning.DeleteDataSourceInput{
		DataSourceId: f.ID,
	})

	return err
}

func (f *MachineLearningDataSource) String() string {
	return *f.ID
}
