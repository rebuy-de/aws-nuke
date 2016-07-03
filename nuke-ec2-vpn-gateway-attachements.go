package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VpnGatewayAttachement struct {
	svc   *ec2.EC2
	vpcId *string
	vpnId *string
}

func (n *EC2Nuke) ListVpnGatewayAttachements() ([]Resource, error) {
	resp, err := n.svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		params := &ec2.DescribeVpnGatewaysInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("attachment.vpc-id"),
					Values: []*string{vpc.VpcId},
				},
			},
		}

		resp, err := n.svc.DescribeVpnGateways(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.VpnGateways {
			resources = append(resources, &EC2VpnGatewayAttachement{
				svc:   n.svc,
				vpcId: vpc.VpcId,
				vpnId: out.VpnGatewayId,
			})
		}
	}

	return resources, nil
}

func (e *EC2VpnGatewayAttachement) Remove() error {
	params := &ec2.DetachVpnGatewayInput{
		VpcId:        e.vpcId,
		VpnGatewayId: e.vpnId,
	}

	_, err := e.svc.DetachVpnGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VpnGatewayAttachement) String() string {
	return fmt.Sprintf("%s->%s", *e.vpnId, *e.vpcId)
}
