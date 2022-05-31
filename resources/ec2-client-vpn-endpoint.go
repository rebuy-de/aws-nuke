package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2ClientVpnEndpoint struct {
	svc     *ec2.EC2
	id      string
	cveTags []*ec2.Tag
}

func init() {
	register("EC2ClientVpnEndpoint", ListEC2ClientVpnEndoint)
}

func ListEC2ClientVpnEndoint(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)
	params := &ec2.DescribeClientVpnEndpointsInput{}

	err := svc.DescribeClientVpnEndpointsPages(params,
		func(page *ec2.DescribeClientVpnEndpointsOutput, lastPage bool) bool {

			for _, out := range page.ClientVpnEndpoints {
				resources = append(resources, &EC2ClientVpnEndpoint{
					svc: svc,
					id:  *out.ClientVpnEndpointId,
				})
			}
			return true
		})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (c *EC2ClientVpnEndpoint) Remove() error {
	params := &ec2.DeleteClientVpnEndpointInput{
		ClientVpnEndpointId: &c.id,
	}

	_, err := c.svc.DeleteClientVpnEndpoint(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2ClientVpnEndpoint) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.cveTags {
		properties.SetTagWithPrefix("cve", tagValue.Key, tagValue.Value)
	}
	return properties
}

func (c *EC2ClientVpnEndpoint) String() string {
	return c.id
}
