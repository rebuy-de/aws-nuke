package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VPNConnection struct {
	svc   *ec2.EC2
	id    string
	state string
}

func (n *EC2Nuke) ListVPNConnections() ([]Resource, error) {
	params := &ec2.DescribeVpnConnectionsInput{}
	resp, err := n.Service.DescribeVpnConnections(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.VpnConnections {
		resources = append(resources, &EC2VPNConnection{
			svc:   n.Service,
			id:    *out.VpnConnectionId,
			state: *out.State,
		})
	}

	return resources, nil
}

func (v *EC2VPNConnection) Filter() error {
	if v.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (v *EC2VPNConnection) Remove() error {
	params := &ec2.DeleteVpnConnectionInput{
		VpnConnectionId: &v.id,
	}

	_, err := v.svc.DeleteVpnConnection(params)
	if err != nil {
		return err
	}

	return nil
}

func (v *EC2VPNConnection) String() string {
	return v.id
}
