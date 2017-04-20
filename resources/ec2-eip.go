package resources

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2Address struct {
	svc *ec2.EC2
	id  string
	ip  string
}

func (n *EC2Nuke) ListAddresses() ([]Resource, error) {
	params := &ec2.DescribeAddressesInput{}
	resp, err := n.Service.DescribeAddresses(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Addresses {
		resources = append(resources, &EC2Address{
			svc: n.Service,
			id:  *out.AllocationId,
			ip:  *out.PublicIp,
		})
	}

	return resources, nil
}

func (e *EC2Address) Remove() error {
	_, err := e.svc.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: &e.id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2Address) String() string {
	return e.ip
}
