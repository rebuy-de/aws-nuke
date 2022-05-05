package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type FSxFileSystem struct {
	svc        *fsx.FSx
	filesystem *fsx.FileSystem
}

func init() {
	register("FSxFileSystem", ListFSxFileSystems)
}

func ListFSxFileSystems(sess *session.Session) ([]Resource, error) {
	svc := fsx.New(sess)
	resources := []Resource{}

	params := &fsx.DescribeFileSystemsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.DescribeFileSystems(params)
		if err != nil {
			return nil, err
		}

		for _, filesystem := range resp.FileSystems {
			resources = append(resources, &FSxFileSystem{
				svc:        svc,
				filesystem: filesystem,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *FSxFileSystem) Remove() error {
	_, err := f.svc.DeleteFileSystem(&fsx.DeleteFileSystemInput{
		FileSystemId: f.filesystem.FileSystemId,
	})

	return err
}

func (f *FSxFileSystem) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range f.filesystem.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("Type", f.filesystem.FileSystemType)
	return properties
}

func (f *FSxFileSystem) String() string {
	return *f.filesystem.FileSystemId
}
