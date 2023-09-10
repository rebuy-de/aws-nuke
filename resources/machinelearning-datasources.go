package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/machinelearning"
	"github.com/sirupsen/logrus"
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
			if aerr, ok := err.(awserr.Error); ok {
				if strings.Contains(aerr.Message(), "AmazonML is no longer available to new customers") {
					logrus.Info("MachineLearningBranchPrediction: AmazonML is no longer available to new customers. Ignore if you haven't set it up.")
					return nil, nil
				}
			}
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
