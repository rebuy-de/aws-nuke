package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2InternetGatewayAttachment struct {
	svc        *ec2.EC2
	vpcId      *string
	vpcTags    []*ec2.Tag
	igwId      *string
	igwTags    []*ec2.Tag
	defaultVPC bool
}

func init() {
	register("EC2InternetGatewayAttachment", ListEC2InternetGatewayAttachments)
}

func ListEC2InternetGatewayAttachments(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		params := &ec2.DescribeInternetGatewaysInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("attachment.vpc-id"),
					Values: []*string{vpc.VpcId},
				},
			},
		}

		resp, err := svc.DescribeInternetGateways(params)
		if err != nil {
			return nil, err
		}

		for _, igw := range resp.InternetGateways {
			resources = append(resources, &EC2InternetGatewayAttachment{
				svc:        svc,
				vpcId:      vpc.VpcId,
				vpcTags:    vpc.Tags,
				igwId:      igw.InternetGatewayId,
				igwTags:    igw.Tags,
				defaultVPC: *vpc.IsDefault,
			})
		}
	}

	return resources, nil
}

func (e *EC2InternetGatewayAttachment) Remove() error {
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

func (e *EC2InternetGatewayAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.igwTags {
		properties.SetTagWithPrefix("igw", tagValue.Key, tagValue.Value)
	}
	for _, tagValue := range e.vpcTags {
		properties.SetTagWithPrefix("vpc", tagValue.Key, tagValue.Value)
	}
	properties.Set("DefaultVPC", e.defaultVPC)
	return properties
}

func (e *EC2InternetGatewayAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *e.igwId, *e.vpcId)
}
