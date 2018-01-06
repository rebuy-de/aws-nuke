package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

type ELB struct {
	svc  *elb.ELB
	name *string
}

func init() {
	register("ElbELB", ListElbELBs)
}

func ListElbELBs(sess *session.Session) ([]Resource, error) {
	svc := elb.New(sess)

	resp, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, elb := range resp.LoadBalancerDescriptions {
		resources = append(resources, &ELB{
			svc:  svc,
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
