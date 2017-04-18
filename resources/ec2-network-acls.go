package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2NetworkACL struct {
	svc       *ec2.EC2
	id        *string
	isDefault *bool
	region    *string
}

func (n *EC2Nuke) ListNetworkACLs() ([]Resource, error) {
	resp, err := n.Service.DescribeNetworkAcls(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.NetworkAcls {

		resources = append(resources, &EC2NetworkACL{
			svc:       n.Service,
			id:        out.NetworkAclId,
			isDefault: out.IsDefault,
			region:    n.Service.Config.Region,
		})
	}

	return resources, nil
}

func (e *EC2NetworkACL) Filter() error {
	if *e.isDefault {
		return fmt.Errorf("cannot delete default VPC")
	}

	return nil
}

func (e *EC2NetworkACL) Remove() error {
	params := &ec2.DeleteNetworkAclInput{
		NetworkAclId: e.id,
	}

	_, err := e.svc.DeleteNetworkAcl(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2NetworkACL) String() string {
	return fmt.Sprintf("%s in %s", *e.id, *e.region)
}
