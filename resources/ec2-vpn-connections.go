package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2VPNConnection struct {
	svc   *ec2.EC2
	conn  *ec2.VpnConnection
}

func init() {
	register("EC2VPNConnection", ListEC2VPNConnections)
}

func ListEC2VPNConnections(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeVpnConnectionsInput{}
	resp, err := svc.DescribeVpnConnections(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.VpnConnections {
		resources = append(resources, &EC2VPNConnection{
			svc:   svc,
			conn:  out,
		})
	}

	return resources, nil
}

func (v *EC2VPNConnection) Filter() error {
	if *v.conn.State == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (v *EC2VPNConnection) Remove() error {
	params := &ec2.DeleteVpnConnectionInput{
		VpnConnectionId: v.conn.VpnConnectionId,
	}

	_, err := v.svc.DeleteVpnConnection(params)
	if err != nil {
		return err
	}

	return nil
}

func (v *EC2VPNConnection) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range v.conn.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (v *EC2VPNConnection) String() string {
	return *v.conn.VpnConnectionId
}
