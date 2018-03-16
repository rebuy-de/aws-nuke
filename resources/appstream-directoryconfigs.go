package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appstream"
)

type AppStreamDirectoryConfig struct {
	svc  *appstream.AppStream
	name *string
}

func init() {
	register("AppStreamDirectoryConfig", ListAppStreamDirectoryConfigs)
}

func ListAppStreamDirectoryConfigs(sess *session.Session) ([]Resource, error) {
	svc := appstream.New(sess)
	resources := []Resource{}

	params := &appstream.DescribeDirectoryConfigsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeDirectoryConfigs(params)
		if err != nil {
			return nil, err
		}

		for _, directoryConfig := range output.DirectoryConfigs {
			resources = append(resources, &AppStreamDirectoryConfig{
				svc:  svc,
				name: directoryConfig.DirectoryName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AppStreamDirectoryConfig) Remove() error {

	_, err := f.svc.DeleteDirectoryConfig(&appstream.DeleteDirectoryConfigInput{
		DirectoryName: f.name,
	})

	return err
}

func (f *AppStreamDirectoryConfig) String() string {
	return *f.name
}
