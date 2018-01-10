package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/efs"
)

type EFSMountTarget struct {
	svc  *efs.EFS
	id   string
	fsid string
}

func init() {
	register("EFSMountTarget", ListEFSMountTargets)
}

func ListEFSMountTargets(sess *session.Session) ([]Resource, error) {
	svc := efs.New(sess)

	resp, err := svc.DescribeFileSystems(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, fs := range resp.FileSystems {
		mt, err := svc.DescribeMountTargets(&efs.DescribeMountTargetsInput{
			FileSystemId: fs.FileSystemId,
		})

		if err != nil {
			return nil, err
		}

		for _, t := range mt.MountTargets {
			resources = append(resources, &EFSMountTarget{
				svc:  svc,
				id:   *t.MountTargetId,
				fsid: *t.FileSystemId,
			})

		}
	}

	return resources, nil
}

func (e *EFSMountTarget) Remove() error {
	_, err := e.svc.DeleteMountTarget(&efs.DeleteMountTargetInput{
		MountTargetId: &e.id,
	})

	return err
}

func (e *EFSMountTarget) String() string {
	return fmt.Sprintf("%s:%s", e.fsid, e.id)
}
