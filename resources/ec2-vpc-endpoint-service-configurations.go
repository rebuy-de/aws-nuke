package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2VPCEndpointServiceConfiguration struct {
	svc  *ec2.EC2
	id   *string
	name *string
}

func init() {
	register("EC2VPCEndpointServiceConfiguration", ListEC2VPCEndpointServiceConfigurations)
}

func ListEC2VPCEndpointServiceConfigurations(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)

	params := &ec2.DescribeVpcEndpointServiceConfigurationsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.DescribeVpcEndpointServiceConfigurations(params)
		if err != nil {
			return nil, err
		}

		for _, serviceConfig := range resp.ServiceConfigurations {
			resources = append(resources, &EC2VPCEndpointServiceConfiguration{
				svc:  svc,
				id:   serviceConfig.ServiceId,
				name: serviceConfig.ServiceName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (e *EC2VPCEndpointServiceConfiguration) Remove() error {
	params := &ec2.DeleteVpcEndpointServiceConfigurationsInput{
		ServiceIds: []*string{e.id},
	}

	_, err := e.svc.DeleteVpcEndpointServiceConfigurations(params)
	if err != nil {
		return err
	}
	return nil
}

func (e *EC2VPCEndpointServiceConfiguration) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", e.name)
	return properties
}

func (e *EC2VPCEndpointServiceConfiguration) String() string {
	return *e.id
}
