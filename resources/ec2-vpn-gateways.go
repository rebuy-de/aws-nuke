package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VpnGateway struct {
	svc   *ec2.EC2
	id    string
	state string
}

func (n *EC2Nuke) ListVpnGateways() ([]Resource, error) {
	params := &ec2.DescribeVpnGatewaysInput{}
	resp, err := n.Service.DescribeVpnGateways(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.VpnGateways {
		resources = append(resources, &EC2VpnGateway{
			svc:   n.Service,
			id:    *out.VpnGatewayId,
			state: *out.State,
		})
	}

	return resources, nil
}

func (i *EC2VpnGateway) Check() error {
	if i.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (e *EC2VpnGateway) Remove() error {
	params := &ec2.DeleteVpnGatewayInput{
		VpnGatewayId: &e.id,
	}

	_, err := e.svc.DeleteVpnGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VpnGateway) String() string {
	return e.id
}
