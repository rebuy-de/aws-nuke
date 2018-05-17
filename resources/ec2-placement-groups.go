package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2PlacementGroup struct {
	svc   *ec2.EC2
	name  string
	state string
}

func init() {
	register("EC2PlacementGroup", ListEC2PlacementGroups)
}

func ListEC2PlacementGroups(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribePlacementGroupsInput{}
	resp, err := svc.DescribePlacementGroups(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.PlacementGroups {
		resources = append(resources, &EC2PlacementGroup{
			svc:   svc,
			name:  *out.GroupName,
			state: *out.State,
		})
	}

	return resources, nil
}

func (p *EC2PlacementGroup) Filter() error {
	if p.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (p *EC2PlacementGroup) Remove() error {
	params := &ec2.DeletePlacementGroupInput{
		GroupName: &p.name,
	}

	_, err := p.svc.DeletePlacementGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (p *EC2PlacementGroup) String() string {
	return p.name
}
