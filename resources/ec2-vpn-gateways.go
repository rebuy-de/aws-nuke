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

func (v *EC2VpnGateway) Filter() error {
	if v.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (v *EC2VpnGateway) Remove() error {
	params := &ec2.DeleteVpnGatewayInput{
		VpnGatewayId: &v.id,
	}

	_, err := v.svc.DeleteVpnGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (v *EC2VpnGateway) String() string {
	return v.id
}
