package main

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2CustomerGateway struct {
	svc *ec2.EC2
	id  *string
}

func (n *EC2Nuke) ListCustomerGateways() ([]Resource, error) {
	params := &ec2.DescribeCustomerGatewaysInput{}
	resp, err := n.svc.DescribeCustomerGateways(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.CustomerGateways {
		resources = append(resources, &EC2CustomerGateway{
			svc: n.svc,
			id:  out.CustomerGatewayId,
		})
	}

	return resources, nil
}

func (e *EC2CustomerGateway) Remove() error {
	params := &ec2.DeleteCustomerGatewayInput{
		CustomerGatewayId: e.id,
	}

	_, err := e.svc.DeleteCustomerGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2CustomerGateway) String() string {
	return *e.id
}
