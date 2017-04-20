package resources

import "github.com/aws/aws-sdk-go/service/cloudtrail"

func (n *CloudTrailNuke) ListTrails() ([]Resource, error) {
	resp, err := n.Service.DescribeTrails(nil)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, trail := range resp.TrailList {
		resources = append(resources, &CloudTrailTrail{
			svc:  n.Service,
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

func (trail *CloudTrailTrail) String() string {
	return *trail.name
}
