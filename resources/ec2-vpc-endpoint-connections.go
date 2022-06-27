package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VPCEndpointConnection struct {
	svc           *ec2.EC2
	serviceID     *string
	vpcEndpointID *string
	state         *string
	owner         *string
}

func init() {
	register("EC2VPCEndpointConnection", ListEC2VPCEndpointConnections)
}

func ListEC2VPCEndpointConnections(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)
	params := &ec2.DescribeVpcEndpointConnectionsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.DescribeVpcEndpointConnections(params)
		if err != nil {
			return nil, err
		}

		for _, endpointConnection := range resp.VpcEndpointConnections {
			resources = append(resources, &EC2VPCEndpointConnection{
				svc:           svc,
				vpcEndpointID: endpointConnection.VpcEndpointId,
				serviceID:     endpointConnection.ServiceId,
				state:         endpointConnection.VpcEndpointState,
				owner:         endpointConnection.VpcEndpointOwner,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (c *EC2VPCEndpointConnection) Filter() error {
	if *c.state == "deleting" || *c.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (c *EC2VPCEndpointConnection) Remove() error {
	params := &ec2.RejectVpcEndpointConnectionsInput{
		ServiceId: c.serviceID,
		VpcEndpointIds: []*string{
			c.vpcEndpointID,
		},
	}

	_, err := c.svc.RejectVpcEndpointConnections(params)
	if err != nil {
		return err
	}
	return nil
}

func (c *EC2VPCEndpointConnection) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("VpcEndpointID", c.vpcEndpointID)
	properties.Set("State", c.state)
	properties.Set("Owner", c.owner)
	return properties
}

func (c *EC2VPCEndpointConnection) String() string {
	return *c.serviceID
}
