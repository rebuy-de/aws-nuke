package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2VPNGatewayAttachment struct {
	svc     *ec2.EC2
	vpcId   string
	vpnId   string
	state   string
	vpcTags []*ec2.Tag
	vgwTags []*ec2.Tag
}

func init() {
	register("EC2VPNGatewayAttachment", ListEC2VPNGatewayAttachments)
}

func ListEC2VPNGatewayAttachments(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		params := &ec2.DescribeVpnGatewaysInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("attachment.vpc-id"),
					Values: []*string{vpc.VpcId},
				},
			},
		}

		resp, err := svc.DescribeVpnGateways(params)
		if err != nil {
			return nil, err
		}

		for _, vgw := range resp.VpnGateways {
			resources = append(resources, &EC2VPNGatewayAttachment{
				svc:     svc,
				vpcId:   *vpc.VpcId,
				vpnId:   *vgw.VpnGatewayId,
				vpcTags: vpc.Tags,
				vgwTags: vgw.Tags,
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

func (v *EC2VPNGatewayAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range v.vgwTags {
		properties.SetTagWithPrefix("vgw", tagValue.Key, tagValue.Value)
	}
	for _, tagValue := range v.vpcTags {
		properties.SetTagWithPrefix("vpc", tagValue.Key, tagValue.Value)
	}
	return properties
}

func (v *EC2VPNGatewayAttachment) String() string {
	return fmt.Sprintf("%s -> %s", v.vpnId, v.vpcId)
}
