package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2Volume struct {
	svc    *ec2.EC2
	volume *ec2.Volume
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
			svc:    svc,
			volume: out,
		})
	}

	return resources, nil
}

func (e *EC2Volume) Remove() error {
	_, err := e.svc.DeleteVolume(&ec2.DeleteVolumeInput{
		VolumeId: e.volume.VolumeId,
	})
	return err
}

func (e *EC2Volume) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("State", e.volume.State)
	for _, tagValue := range e.volume.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *EC2Volume) String() string {
	return *e.volume.VolumeId
}
