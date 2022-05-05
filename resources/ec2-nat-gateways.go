package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2NATGateway struct {
	svc   *ec2.EC2
	natgw *ec2.NatGateway
}

func init() {
	register("EC2NATGateway", ListEC2NATGateways)
}

func ListEC2NATGateways(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeNatGatewaysInput{}
	resp, err := svc.DescribeNatGateways(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, natgw := range resp.NatGateways {
		resources = append(resources, &EC2NATGateway{
			svc:   svc,
			natgw: natgw,
		})
	}

	return resources, nil
}

func (n *EC2NATGateway) Filter() error {
	if *n.natgw.State == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (n *EC2NATGateway) Remove() error {
	params := &ec2.DeleteNatGatewayInput{
		NatGatewayId: n.natgw.NatGatewayId,
	}

	_, err := n.svc.DeleteNatGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (n *EC2NATGateway) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range n.natgw.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (n *EC2NATGateway) String() string {
	return *n.natgw.NatGatewayId
}
