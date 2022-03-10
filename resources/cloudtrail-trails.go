package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("CloudTrailTrail", ListCloudTrailTrails)
}

func ListCloudTrailTrails(sess *session.Session) ([]Resource, error) {
	svc := cloudtrail.New(sess)

	resp, err := svc.DescribeTrails(nil)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, trail := range resp.TrailList {
		resources = append(resources, &CloudTrailTrail{
			svc:  svc,
			name: trail.Name,
		})

	}
	return resources, nil
}

type CloudTrailTrail struct {
	svc  *cloudtrail.CloudTrail
	name *string
}

func (trail *CloudTrailTrail) Remove() error {
	_, err := trail.svc.DeleteTrail(&cloudtrail.DeleteTrailInput{
		Name: trail.name,
	})
	return err
}

func (trail *CloudTrailTrail) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", trail.name)
}

func (trail *CloudTrailTrail) String() string {
	return *trail.name
}
