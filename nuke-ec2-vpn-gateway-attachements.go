package main

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2VpnGatewayAttachement struct {
	svc *ec2.EC2
	id  *string
}

func (n *EC2Nuke) ListVpnGatewayAttachements() ([]Resource, error) {
	params := &ec2.DescribeVpnGatewaysInput{}
	resp, err := n.svc.DescribeVpnGateways(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.VpnGateways {
		resources = append(resources, &EC2VpnGatewayAttachement{
			svc: n.svc,
			id:  out.VpnGatewayId,
		})
	}

	return resources, nil
}

func (e *EC2VpnGatewayAttachement) Remove() error {
	params := &ec2.DetachVpnGatewayInput{
		VpnGatewayId: e.id,
	}

	_, err := e.svc.DetachVpnGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VpnGatewayAttachement) String() string {
	return *e.id
}
