package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2KeyPair struct {
	svc  *ec2.EC2
	name string
	tags []*ec2.Tag
}

func init() {
	register("EC2KeyPair", ListEC2KeyPairs)
}

func ListEC2KeyPairs(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeKeyPairs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.KeyPairs {
		resources = append(resources, &EC2KeyPair{
			svc:  svc,
			name: *out.KeyName,
			tags: out.Tags,
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

func (e *EC2KeyPair) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", e.name)

	for _, tag := range e.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (e *EC2KeyPair) String() string {
	return e.name
}
