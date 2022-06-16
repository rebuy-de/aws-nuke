package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2SecurityGroup struct {
	svc     *ec2.EC2
	group   *ec2.SecurityGroup
	id      *string
	name    *string
	ingress []*ec2.IpPermission
	egress  []*ec2.IpPermission
	ownerID *string
}

func init() {
	register("EC2SecurityGroup", ListEC2SecurityGroups)
}

func ListEC2SecurityGroups(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)

	params := &ec2.DescribeSecurityGroupsInput{}
	err := svc.DescribeSecurityGroupsPages(params,
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			for _, group := range page.SecurityGroups {
				resources = append(resources, &EC2SecurityGroup{
					svc:     svc,
					group:   group,
					id:      group.GroupId,
					name:    group.GroupName,
					ingress: group.IpPermissions,
					egress:  group.IpPermissionsEgress,
					ownerID: group.OwnerId,
				})
			}
			return !lastPage
		})

	if err != nil {
		return nil, err
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

func (sg *EC2SecurityGroup) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range sg.group.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("Name", sg.name)
	properties.Set("OwnerID", sg.ownerID)
	return properties
}

func (sg *EC2SecurityGroup) String() string {
	return *sg.id
}
