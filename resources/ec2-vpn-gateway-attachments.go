package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VPNGatewayAttachment struct {
	svc   *ec2.EC2
	vpcId string
	vpnId string
	state string
}

func init() {
	register("EC2VPNGatewayAttachment", ListEC2VPNGatewayAttachments)
}

func ListEC2VPNGatewayAttachments(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeVpnGateways(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, vgw := range resp.VpnGateways {
		for _, att := range vgw.VpcAttachments {
			resources = append(resources, &EC2VPNGatewayAttachment{
				svc:   svc,
				vpcId: *att.VpcId,
				vpnId: *vgw.VpnGatewayId,
				state: *att.State,
			})
		}
	}

	return resources, nil
}

func (v *EC2VPNGatewayAttachment) Filter() error {
	if v.state == "detached" {
		return fmt.Errorf("already detached")
	}
	return nil
}

func (v *EC2VPNGatewayAttachment) Remove() error {
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

func (v *EC2VPNGatewayAttachment) String() string {
	return fmt.Sprintf("%s -> %s", v.vpnId, v.vpcId)
}
