package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2SecurityGroup struct {
	svc  *ec2.EC2
	id   *string
	name *string
}

func (n *EC2Nuke) ListSecurityGroups() ([]Resource, error) {
	params := &ec2.DescribeSecurityGroupsInput{}
	resp, err := n.Service.DescribeSecurityGroups(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, group := range resp.SecurityGroups {
		resources = append(resources, &EC2SecurityGroup{
			svc:  n.Service,
			id:   group.GroupId,
			name: group.GroupName,
		})
	}

	return resources, nil
}

func (sg *EC2SecurityGroup) Filter() error {
	if *sg.name == "default" {
		return fmt.Errorf("cannot delete group 'default'")
	}

	return nil
}

func (sg *EC2SecurityGroup) Remove() error {
	params := &ec2.DeleteSecurityGroupInput{
		GroupId: sg.id,
	}

	_, err := sg.svc.DeleteSecurityGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (sg *EC2SecurityGroup) String() string {
	return *sg.id
}
