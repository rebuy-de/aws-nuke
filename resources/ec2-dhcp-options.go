package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/Optum/aws-nuke/pkg/types"
)

type EC2DHCPOption struct {
	svc  *ec2.EC2
	id   *string
	tags []*ec2.Tag
}

func init() {
	register("EC2DHCPOption", ListEC2DHCPOptions)
}

func ListEC2DHCPOptions(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeDhcpOptions(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.DhcpOptions {

		resources = append(resources, &EC2DHCPOption{
			svc:  svc,
			id:   out.DhcpOptionsId,
			tags: out.Tags,
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

func (e *EC2DHCPOption) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *EC2DHCPOption) String() string {
	return *e.id
}
