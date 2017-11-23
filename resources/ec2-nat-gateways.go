package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2NatGateway struct {
	svc   *ec2.EC2
	id    string
	state string
}

func (n *EC2Nuke) ListNatGateways() ([]Resource, error) {
	params := &ec2.DescribeNatGatewaysInput{}
	resp, err := n.Service.DescribeNatGateways(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.NatGateways {
		resources = append(resources, &EC2NatGateway{
			svc:   n.Service,
			id:    *out.NatGatewayId,
			state: *out.State,
		})
	}

	return resources, nil
}

func (n *EC2NatGateway) Filter() error {
	if n.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (n *EC2NatGateway) Remove() error {
	params := &ec2.DeleteNatGatewayInput{
		NatGatewayId: &n.id,
	}

	_, err := n.svc.DeleteNatGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (n *EC2NatGateway) String() string {
	return n.id
}
