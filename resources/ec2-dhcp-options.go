package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2DhcpOption struct {
	svc    *ec2.EC2
	id     *string
	region *string
}

func (n *EC2Nuke) ListDhcpOptions() ([]Resource, error) {
	resp, err := n.Service.DescribeDhcpOptions(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.DhcpOptions {

		resources = append(resources, &EC2DhcpOption{
			svc:    n.Service,
			id:     out.DhcpOptionsId,
			region: n.Service.Config.Region,
		})
	}

	return resources, nil
}

func (e *EC2DhcpOption) Remove() error {
	params := &ec2.DeleteDhcpOptionsInput{
		DhcpOptionsId: e.id,
	}

	_, err := e.svc.DeleteDhcpOptions(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2DhcpOption) String() string {
	return fmt.Sprintf("%s in %s", *e.id, *e.region)
}
