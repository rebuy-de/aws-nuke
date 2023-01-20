package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/aws/aws-sdk-go/aws"
)

type EC2VPCEndpoint struct {
	svc     *ec2.EC2
	id      *string
	vpcTags []*ec2.Tag
}

func init() {
	register("EC2VPCEndpoint", ListEC2VPCEndpoints)
}

func ListEC2VPCEndpoints(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		params := &ec2.DescribeVpcEndpointsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{vpc.VpcId},
				},
			},
		}

		resp, err := svc.DescribeVpcEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, vpcEndpoint := range resp.VpcEndpoints {
			resources = append(resources, &EC2VPCEndpoint{
				svc:  svc,
				id:   vpcEndpoint.VpcEndpointId,
				vpcTags: vpc.Tags,
			})
		}
	}

	return resources, nil
}

func (endpoint *EC2VPCEndpoint) Remove() error {
	params := &ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []*string{endpoint.id},
	}

	_, err := endpoint.svc.DeleteVpcEndpoints(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VPCEndpoint) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.vpcTags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (endpoint *EC2VPCEndpoint) String() string {
	return *endpoint.id
}
