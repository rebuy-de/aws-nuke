package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Vpc struct {
	svc    *ec2.EC2
	id     *string
	region *string
}

func (n *EC2Nuke) ListVpcs() ([]Resource, error) {
	resp, err := n.Service.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		resources = append(resources, &EC2Vpc{
			svc:    n.Service,
			id:     vpc.VpcId,
			region: n.Service.Config.Region,
		})
	}

	return resources, nil
}

func (e *EC2Vpc) Remove() error {
	params := &ec2.DeleteVpcInput{
		VpcId: e.id,
	}

	_, err := e.svc.DeleteVpc(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2Vpc) String() string {
	return fmt.Sprintf("%s in %s", *e.id, *e.region)
}
