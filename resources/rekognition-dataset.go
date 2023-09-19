package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type RekognitionDataset struct {
	svc *rekognition.Rekognition
	arn *string
}

func init() {
	register("RekognitionDataset", ListRekognitionDatasets)
}

func ListRekognitionDatasets(sess *session.Session) ([]Resource, error) {
	svc := rekognition.New(sess)
	resources := []Resource{}

	params := &rekognition.DescribeProjectsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeProjects(params)
		if err != nil {
			return nil, err
		}

		for _, project := range output.ProjectDescriptions {
			for _, dataset := range project.Datasets {
				resources = append(resources, &RekognitionDataset{
					svc: svc,
					arn: dataset.DatasetArn,
				})
			}
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *RekognitionDataset) Remove() error {

	_, err := f.svc.DeleteDataset(&rekognition.DeleteDatasetInput{
		DatasetArn: f.arn,
	})

	return err
}

func (f *RekognitionDataset) String() string {
	return *f.arn
}
