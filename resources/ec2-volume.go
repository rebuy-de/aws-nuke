package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Volume struct {
	svc *ec2.EC2
	id  string
}

func init() {
	register("EC2Volume", ListEC2Volumes)
}

func ListEC2Volumes(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeVolumes(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Volumes {
		resources = append(resources, &EC2Volume{
			svc: svc,
			id:  *out.VolumeId,
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
	return e.id
}
