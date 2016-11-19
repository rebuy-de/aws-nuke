package resources

import "github.com/aws/aws-sdk-go/service/elb"

type ELB struct {
	svc  *elb.ELB
	name *string
}

func (n *ElbNuke) ListELBs() ([]Resource, error) {
	resp, err := n.Service.DescribeLoadBalancers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, elb := range resp.LoadBalancerDescriptions {
		resources = append(resources, &ELB{
			svc:  n.Service,
			name: elb.LoadBalancerName,
		})
	}

	return resources, nil
}

func (e *ELB) Remove() error {
	params := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: e.name,
	}

	_, err := e.svc.DeleteLoadBalancer(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *ELB) String() string {
	return *e.name
}
