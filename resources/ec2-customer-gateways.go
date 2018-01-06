package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2CustomerGateway struct {
	svc   *ec2.EC2
	id    string
	state string
}

func init() {
	register("EC2CustomerGateway", ListEC2CustomerGateways)
}

func ListEC2CustomerGateways(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeCustomerGatewaysInput{}
	resp, err := svc.DescribeCustomerGateways(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.CustomerGateways {
		resources = append(resources, &EC2CustomerGateway{
			svc:   svc,
			id:    *out.CustomerGatewayId,
			state: *out.State,
		})
	}

	return resources, nil
}

func (c *EC2CustomerGateway) Filter() error {
	if c.state == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (c *EC2CustomerGateway) Remove() error {
	params := &ec2.DeleteCustomerGatewayInput{
		CustomerGatewayId: &c.id,
	}

	_, err := c.svc.DeleteCustomerGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (c *EC2CustomerGateway) String() string {
	return c.id
}
