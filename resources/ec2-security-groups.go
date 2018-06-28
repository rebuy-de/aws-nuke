package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2SecurityGroup struct {
	svc     *ec2.EC2
	id      *string
	name    *string
	ingress []*ec2.IpPermission
	egress  []*ec2.IpPermission
}

func init() {
	register("EC2SecurityGroup", ListEC2SecurityGroups)
}

func ListEC2SecurityGroups(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeSecurityGroupsInput{}
	resp, err := svc.DescribeSecurityGroups(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, group := range resp.SecurityGroups {
		resources = append(resources, &EC2SecurityGroup{
			svc:     svc,
			id:      group.GroupId,
			name:    group.GroupName,
			ingress: group.IpPermissions,
			egress:  group.IpPermissionsEgress,
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
	if len(sg.egress) > 0 {
		egressParams := &ec2.RevokeSecurityGroupEgressInput{
			GroupId:       sg.id,
			IpPermissions: sg.egress,
		}

		_, _ = sg.svc.RevokeSecurityGroupEgress(egressParams)
	}

	if len(sg.ingress) > 0 {
		ingressParams := &ec2.RevokeSecurityGroupIngressInput{
			GroupId:       sg.id,
			IpPermissions: sg.ingress,
		}

		_, _ = sg.svc.RevokeSecurityGroupIngress(ingressParams)
	}

	params := &ec2.DeleteSecurityGroupInput{
		GroupId: sg.id,
	}

	_, err := sg.svc.DeleteSecurityGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (sg *EC2SecurityGroup) Properties() Properties {
	return NewProperties().
		Set("Name", sg.name)
}

func (sg *EC2SecurityGroup) String() string {
	return *sg.id
}
