package resources

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2DHCPOption struct {
	svc *ec2.EC2
	id  *string
}

func (n *EC2Nuke) ListDHCPOptions() ([]Resource, error) {
	resp, err := n.Service.DescribeDhcpOptions(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.DhcpOptions {

		resources = append(resources, &EC2DHCPOption{
			svc: n.Service,
			id:  out.DhcpOptionsId,
		})
	}

	return resources, nil
}

func (e *EC2DHCPOption) Remove() error {
	params := &ec2.DeleteDhcpOptionsInput{
		DhcpOptionsId: e.id,
	}

	_, err := e.svc.DeleteDhcpOptions(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2DHCPOption) String() string {
	return *e.id
}
