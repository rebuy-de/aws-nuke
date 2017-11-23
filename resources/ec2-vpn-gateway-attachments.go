package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VpnGatewayAttachment struct {
	svc   *ec2.EC2
	vpcId string
	vpnId string
	state string
}

func (n *EC2Nuke) ListVpnGatewayAttachments() ([]Resource, error) {

	resp, err := n.Service.DescribeVpnGateways(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, vgw := range resp.VpnGateways {
		for _, att := range vgw.VpcAttachments {
			resources = append(resources, &EC2VpnGatewayAttachment{
				svc:   n.Service,
				vpcId: *att.VpcId,
				vpnId: *vgw.VpnGatewayId,
				state: *att.State,
			})
		}
	}

	return resources, nil
}

func (v *EC2VpnGatewayAttachment) Filter() error {
	if v.state == "detached" {
		return fmt.Errorf("already detached")
	}
	return nil
}

func (v *EC2VpnGatewayAttachment) Remove() error {
	params := &ec2.DetachVpnGatewayInput{
		VpcId:        &v.vpcId,
		VpnGatewayId: &v.vpnId,
	}

	_, err := v.svc.DetachVpnGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (v *EC2VpnGatewayAttachment) String() string {
	return fmt.Sprintf("%s -> %s", v.vpnId, v.vpcId)
}
