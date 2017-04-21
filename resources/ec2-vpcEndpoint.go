package resources

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2VpcEndpoint struct {
	svc *ec2.EC2
	id  *string
}

func (n *EC2Nuke) ListVpcEndpoints() ([]Resource, error) {
	resp, err := n.Service.DescribeVpcEndpoints(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpcEndpoint := range resp.VpcEndpoints {
		resources = append(resources, &EC2VpcEndpoint{
			svc: n.Service,
			id:  vpcEndpoint.VpcEndpointId,
		})
	}

	return resources, nil
}

func (endpoint *EC2VpcEndpoint) Remove() error {
	params := &ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []*string{endpoint.id},
	}

	_, err := endpoint.svc.DeleteVpcEndpoints(params)
	if err != nil {
		return err
	}

	return nil
}

func (endpoint *EC2VpcEndpoint) String() string {
	return *endpoint.id
}
