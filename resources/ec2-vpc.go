package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2VPC struct {
	svc *ec2.EC2
	id  *string
}

func init() {
	register("EC2VPC", ListEC2VPCs)
}

func ListEC2VPCs(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		resources = append(resources, &EC2VPC{
			svc: svc,
			id:  vpc.VpcId,
		})
	}

	return resources, nil
}

func (e *EC2VPC) Remove() error {
	params := &ec2.DeleteVpcInput{
		VpcId: e.id,
	}

	_, err := e.svc.DeleteVpc(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VPC) String() string {
	return *e.id
}
