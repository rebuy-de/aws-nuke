package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2NetworkACL struct {
	svc       *ec2.EC2
	id        *string
	isDefault *bool
	tags      []*ec2.Tag
	ownerID   *string
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
			tags:      out.Tags,
			ownerID:   out.OwnerId,
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

func (f *EC2NetworkACL) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.Set("ID", f.id)
	properties.Set("OwnerID", f.ownerID)
	return properties
}

func (e *EC2NetworkACL) String() string {
	return *e.id
}
