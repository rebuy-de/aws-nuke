package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/clouddirectory"
)

type CloudDirectoryDirectory struct {
	svc          *clouddirectory.CloudDirectory
	directoryARN *string
}

func init() {
	register("CloudDirectoryDirectory", ListCloudDirectoryDirectories)
}

func ListCloudDirectoryDirectories(sess *session.Session) ([]Resource, error) {
	svc := clouddirectory.New(sess)
	resources := []Resource{}

	params := &clouddirectory.ListDirectoriesInput{
		MaxResults: aws.Int64(30),
		State:      aws.String("ENABLED"),
	}

	for {
		resp, err := svc.ListDirectories(params)
		if err != nil {
			return nil, err
		}

		for _, directory := range resp.Directories {
			resources = append(resources, &CloudDirectoryDirectory{
				svc:          svc,
				directoryARN: directory.DirectoryArn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CloudDirectoryDirectory) Remove() error {

	_, err := f.svc.DisableDirectory(&clouddirectory.DisableDirectoryInput{
		DirectoryArn: f.directoryARN,
	})

	if err == nil {
		_, err = f.svc.DeleteDirectory(&clouddirectory.DeleteDirectoryInput{
			DirectoryArn: f.directoryARN,
		})
	}

	return err
}

func (f *CloudDirectoryDirectory) String() string {
	return *f.directoryARN
}
