package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appstream"
)

type AppStreamImageBuilder struct {
	svc  *appstream.AppStream
	name *string
}

func init() {
	register("AppStreamImageBuilder", ListAppStreamImageBuilders)
}

func ListAppStreamImageBuilders(sess *session.Session) ([]Resource, error) {
	svc := appstream.New(sess)
	resources := []Resource{}

	params := &appstream.DescribeImageBuildersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeImageBuilders(params)
		if err != nil {
			return nil, err
		}

		for _, imageBuilder := range output.ImageBuilders {
			resources = append(resources, &AppStreamImageBuilder{
				svc:  svc,
				name: imageBuilder.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AppStreamImageBuilder) Remove() error {

	_, err := f.svc.DeleteImageBuilder(&appstream.DeleteImageBuilderInput{
		Name: f.name,
	})

	return err
}

func (f *AppStreamImageBuilder) String() string {
	return *f.name
}
