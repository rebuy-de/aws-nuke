package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2KeyPair struct {
	svc    *ec2.EC2
	name   string
	region string
}

func (n *EC2Nuke) ListKeyPairs() ([]Resource, error) {
	resp, err := n.Service.DescribeKeyPairs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.KeyPairs {
		resources = append(resources, &EC2KeyPair{
			svc:    n.Service,
			name:   *out.KeyName,
			region: *n.Service.Config.Region,
		})
	}

	return resources, nil
}

func (e *EC2KeyPair) Remove() error {
	params := &ec2.DeleteKeyPairInput{
		KeyName: &e.name,
	}

	_, err := e.svc.DeleteKeyPair(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2KeyPair) String() string {
	return fmt.Sprintf("%s in %s", e.name, e.region)
}
