package main

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2InternetGateways struct {
	svc *ec2.EC2
	id  *string
}

func (n *EC2Nuke) ListInternetGateways() ([]Resource, error) {
	resp, err := n.svc.DescribeInternetGateways(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.InternetGateways {
		resources = append(resources, &EC2InternetGateways{
			svc: n.svc,
			id:  out.InternetGatewayId,
		})
	}

	return resources, nil
}

func (e *EC2InternetGateways) Remove() error {
	params := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: e.id,
	}

	_, err := e.svc.DeleteInternetGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2InternetGateways) String() string {
	return *e.id
}
