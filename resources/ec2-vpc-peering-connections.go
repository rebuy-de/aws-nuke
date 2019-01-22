package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VPCPeeringConnection struct {
	svc *ec2.EC2
	id  *string
}

func init() {
	register("EC2VPCPeeringConnection", ListEC2VPCPeeringConnections)
}

func ListEC2VPCPeeringConnections(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)

	// filter should be set as deleted vpc connetions are returned
	params := &ec2.DescribeVpcPeeringConnectionsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status-code"),
				Values: []*string{aws.String("pending-acceptance"), aws.String("failed"), aws.String("expired"), aws.String("provisioning"), aws.String("active")},
			},
		},
	}

	resp, err := svc.DescribeVpcPeeringConnections(params)
	if err != nil {
		return nil, err
	}

	for _, peeringConfig := range resp.VpcPeeringConnections {
		resources = append(resources, &EC2VPCPeeringConnection{
			svc: svc,
			id:  peeringConfig.VpcPeeringConnectionId,
		})
	}

	return resources, nil
}

func (e *EC2VPCPeeringConnection) Remove() error {
	params := &ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: e.id,
	}

	_, err := e.svc.DeleteVpcPeeringConnection(params)
	if err != nil {
		return err
	}
	return nil
}

func (e *EC2VPCPeeringConnection) String() string {
	return *e.id
}
