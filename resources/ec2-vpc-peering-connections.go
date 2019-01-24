package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VPCPeeringConnection struct {
	svc    *ec2.EC2
	id     *string
	status *string
}

func init() {
	register("EC2VPCPeeringConnection", ListEC2VPCPeeringConnections)
}

func ListEC2VPCPeeringConnections(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)

	// filter should be set as deleted vpc connetions are returned
	params := &ec2.DescribeVpcPeeringConnectionsInput{}

	resp, err := svc.DescribeVpcPeeringConnections(params)
	if err != nil {
		return nil, err
	}

	for _, peeringConfig := range resp.VpcPeeringConnections {
		resources = append(resources, &EC2VPCPeeringConnection{
			svc:    svc,
			id:     peeringConfig.VpcPeeringConnectionId,
			status: peeringConfig.Status.Code,
		})
	}

	return resources, nil
}

func (p *EC2VPCPeeringConnection) Filter() error {
	if *p.status == "deleting" || *p.status == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (p *EC2VPCPeeringConnection) Remove() error {
	params := &ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: p.id,
	}

	_, err := p.svc.DeleteVpcPeeringConnection(params)
	if err != nil {
		return err
	}
	return nil
}

func (p *EC2VPCPeeringConnection) String() string {
	return *p.id
}
