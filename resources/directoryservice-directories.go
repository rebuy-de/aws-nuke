package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/directoryservice"
)

type DirectoryServiceDirectory struct {
	svc         *directoryservice.DirectoryService
	directoryID *string
}

func init() {
	register("DirectoryServiceDirectory", ListDirectoryServiceDirectories)
}

func ListDirectoryServiceDirectories(sess *session.Session) ([]Resource, error) {
	svc := directoryservice.New(sess)
	resources := []Resource{}

	params := &directoryservice.DescribeDirectoriesInput{
		Limit: aws.Int64(100),
	}

	for {
		resp, err := svc.DescribeDirectories(params)
		if err != nil {
			return nil, err
		}

		for _, directory := range resp.DirectoryDescriptions {
			resources = append(resources, &DirectoryServiceDirectory{
				svc:         svc,
				directoryID: directory.DirectoryId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *DirectoryServiceDirectory) Remove() error {

	_, err := f.svc.DeleteDirectory(&directoryservice.DeleteDirectoryInput{
		DirectoryId: f.directoryID,
	})

	return err
}

func (f *DirectoryServiceDirectory) String() string {
	return *f.directoryID
}
