package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2NetworkACL struct {
	svc       *ec2.EC2
	id        *string
	isDefault *bool
}

func init() {
	register("EC2NetworkACL", ListEC2NetworkACLs)
}

func ListEC2NetworkACLs(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeNetworkAcls(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.NetworkAcls {

		resources = append(resources, &EC2NetworkACL{
			svc:       svc,
			id:        out.NetworkAclId,
			isDefault: out.IsDefault,
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
	return *e.id
}
