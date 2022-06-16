package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EFSFileSystem struct {
	svc     *efs.EFS
	id      string
	name    string
	tagList []*efs.Tag
}

func init() {
	register("EFSFileSystem", ListEFSFileSystems)
}

func ListEFSFileSystems(sess *session.Session) ([]Resource, error) {
	svc := efs.New(sess)

	resp, err := svc.DescribeFileSystems(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, fs := range resp.FileSystems {
		lto, err := svc.ListTagsForResource(&efs.ListTagsForResourceInput{ResourceId: fs.FileSystemId})
		if err != nil {
			return nil, err
		}
		resources = append(resources, &EFSFileSystem{
			svc:     svc,
			id:      *fs.FileSystemId,
			name:    *fs.CreationToken,
			tagList: lto.Tags,
		})

	}

	return resources, nil
}

func (e *EFSFileSystem) Remove() error {
	_, err := e.svc.DeleteFileSystem(&efs.DeleteFileSystemInput{
		FileSystemId: &e.id,
	})

	return err
}

func (e *EFSFileSystem) Properties() types.Properties {
	properties := types.NewProperties()
	for _, t := range e.tagList {
		properties.SetTag(t.Key, t.Value)
	}
	return properties
}

func (e *EFSFileSystem) String() string {
	return e.name
}
