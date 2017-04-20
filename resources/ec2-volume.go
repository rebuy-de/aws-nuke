package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Volume struct {
	svc    *ec2.EC2
	id     string
	region string
}

func (n *EC2Nuke) ListVolumes() ([]Resource, error) {
	resp, err := n.Service.DescribeVolumes(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Volumes {
		resources = append(resources, &EC2Volume{
			svc:    n.Service,
			id:     *out.VolumeId,
			region: *n.Service.Config.Region,
		})
	}

	return resources, nil
}

func (e *EC2Volume) Remove() error {
	_, err := e.svc.DeleteVolume(&ec2.DeleteVolumeInput{
		VolumeId: &e.id,
	})
	return err
}

func (e *EC2Volume) String() string {
	return fmt.Sprintf("%s in %s", e.id, e.region)
}
