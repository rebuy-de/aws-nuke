package resources

import "github.com/aws/aws-sdk-go/service/efs"

type EFSFileSystem struct {
	svc  *efs.EFS
	id   string
	name string
}

func (n *EFSNuke) ListFileSystems() ([]Resource, error) {
	resp, err := n.Service.DescribeFileSystems(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, fs := range resp.FileSystems {
		resources = append(resources, &EFSFileSystem{
			svc:  n.Service,
			id:   *fs.FileSystemId,
			name: *fs.CreationToken,
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

func (e *EFSFileSystem) String() string {
	return e.name
}
