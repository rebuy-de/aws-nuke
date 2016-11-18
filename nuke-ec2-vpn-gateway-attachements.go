package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VpnGatewayAttachement struct {
	svc   *ec2.EC2
	vpcId string
	vpnId string
	state string
}

func (n *EC2Nuke) ListVpnGatewayAttachements() ([]Resource, error) {

	resp, err := n.svc.DescribeVpnGateways(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, vgw := range resp.VpnGateways {
		for _, att := range vgw.VpcAttachments {
			resources = append(resources, &EC2VpnGatewayAttachement{
				svc:   n.svc,
				vpcId: *att.VpcId,
				vpnId: *vgw.VpnGatewayId,
				state: *att.State,
			})
		}
	}

	return resources, nil
}

func (i *EC2VpnGatewayAttachement) Check() error {
	if i.state == "detached" {
		return fmt.Errorf("already detached")
	}
	return nil
}

func (e *EC2VpnGatewayAttachement) Remove() error {
	params := &ec2.DetachVpnGatewayInput{
		VpcId:        &e.vpcId,
		VpnGatewayId: &e.vpnId,
	}

	_, err := e.svc.DetachVpnGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VpnGatewayAttachement) String() string {
	return fmt.Sprintf("%s -> %s", e.vpnId, e.vpcId)
}
