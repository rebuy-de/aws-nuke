package main

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2VpnConnection struct {
	svc *ec2.EC2
	id  *string
}

func (n *EC2Nuke) ListVpnConnections() ([]Resource, error) {
	params := &ec2.DescribeVpnConnectionsInput{}
	resp, err := n.svc.DescribeVpnConnections(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.VpnConnections {
		resources = append(resources, &EC2VpnConnection{
			svc: n.svc,
			id:  out.VpnConnectionId,
		})
	}

	return resources, nil
}

func (e *EC2VpnConnection) Remove() error {
	params := &ec2.DeleteVpnConnectionInput{
		VpnConnectionId: e.id,
	}

	_, err := e.svc.DeleteVpnConnection(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VpnConnection) String() string {
	return *e.id
}
