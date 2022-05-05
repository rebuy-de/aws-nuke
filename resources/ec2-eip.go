package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2Address struct {
	svc *ec2.EC2
	eip *ec2.Address
	id  string
	ip  string
}

func init() {
	register("EC2Address", ListEC2Addresses)
}

func ListEC2Addresses(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeAddressesInput{}
	resp, err := svc.DescribeAddresses(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Addresses {
		resources = append(resources, &EC2Address{
			svc: svc,
			eip: out,
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

func (e *EC2Address) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.eip.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("AllocationID", e.id)
	return properties
}

func (e *EC2Address) String() string {
	return e.ip
}
