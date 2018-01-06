package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VPNGateway struct {
	svc   *ec2.EC2
	id    string
	state string
}

func init() {
	register("EC2VPNGateway", ListEC2VPNGateways)
}

func ListEC2VPNGateways(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeVpnGatewaysInput{}
	resp, err := svc.DescribeVpnGateways(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.VpnGateways {
		resources = append(resources, &EC2VPNGateway{
			svc:   svc,
			id:    *out.VpnGatewayId,
			state: *out.State,
		})
	}

	return resources, nil
}

func (v *EC2VPNGateway) Filter() error {
	if v.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (v *EC2VPNGateway) Remove() error {
	params := &ec2.DeleteVpnGatewayInput{
		VpnGatewayId: &v.id,
	}

	_, err := v.svc.DeleteVpnGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (v *EC2VPNGateway) String() string {
	return v.id
}
