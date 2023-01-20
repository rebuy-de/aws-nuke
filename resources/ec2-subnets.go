package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2Subnet struct {
	svc        *ec2.EC2
	subnet     *ec2.Subnet
	defaultVPC bool
}

func init() {
	register("EC2Subnet", ListEC2Subnets)
}

func ListEC2Subnets(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeSubnetsInput{}
	resp, err := svc.DescribeSubnets(params)
	if err != nil {
		return nil, err
	}

	defVpcId := ""
	if defVpc := DefaultVpc(svc); defVpc != nil {
		defVpcId = *defVpc.VpcId
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Subnets {
		resources = append(resources, &EC2Subnet{
			svc:        svc,
			subnet:     out,
			defaultVPC: defVpcId == *out.VpcId,
		})
	}

	return resources, nil
}

func (e *EC2Subnet) Remove() error {
	params := &ec2.DeleteSubnetInput{
		SubnetId: e.subnet.SubnetId,
	}

	_, err := e.svc.DeleteSubnet(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2Subnet) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.subnet.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("DefaultForAz", e.subnet.DefaultForAz)
	properties.Set("DefaultVPC", e.defaultVPC)
	properties.Set("OwnerID", e.subnet.OwnerId)
	return properties
}

func (e *EC2Subnet) String() string {
	return *e.subnet.SubnetId
}
