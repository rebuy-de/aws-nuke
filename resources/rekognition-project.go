package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type RekognitionProject struct {
	svc *rekognition.Rekognition
	arn *string
}

func init() {
	register("RekognitionProject", ListRekognitionProjects)
}

func ListRekognitionProjects(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &RekognitionProject{
				svc: svc,
				arn: project.ProjectArn,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *RekognitionProject) Remove() error {

	_, err := f.svc.DeleteProject(&rekognition.DeleteProjectInput{
		ProjectArn: f.arn,
	})

	return err
}

func (f *RekognitionProject) String() string {
	return *f.arn
}
