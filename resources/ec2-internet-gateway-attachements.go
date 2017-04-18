package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2InternetGatewayAttachement struct {
	svc    *ec2.EC2
	vpcId  *string
	igwId  *string
	region *string
}

func (n *EC2Nuke) ListInternetGatewayAttachements() ([]Resource, error) {
	resp, err := n.Service.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		params := &ec2.DescribeInternetGatewaysInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("attachment.vpc-id"),
					Values: []*string{vpc.VpcId},
				},
			},
		}

		resp, err := n.Service.DescribeInternetGateways(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.InternetGateways {
			resources = append(resources, &EC2InternetGatewayAttachement{
				svc:    n.Service,
				vpcId:  vpc.VpcId,
				igwId:  out.InternetGatewayId,
				region: n.Service.Config.Region,
			})
		}
	}

	return resources, nil
}

func (e *EC2InternetGatewayAttachement) Remove() error {
	params := &ec2.DetachInternetGatewayInput{
		VpcId:             e.vpcId,
		InternetGatewayId: e.igwId,
	}

	_, err := e.svc.DetachInternetGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2InternetGatewayAttachement) String() string {
	return fmt.Sprintf("%s -> %s in %s", *e.igwId, *e.vpcId, *e.region)
}
