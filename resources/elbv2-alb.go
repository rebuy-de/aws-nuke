package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type ELBv2LoadBalancer struct {
	svc  *elbv2.ELBV2
	name *string
	arn  *string
}

func init() {
	register("ELBv2", ListELBv2LoadBalancers)
}

func ListELBv2LoadBalancers(sess *session.Session) ([]Resource, error) {
	svc := elbv2.New(sess)

	resp, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, elbv2 := range resp.LoadBalancers {
		resources = append(resources, &ELBv2LoadBalancer{
			svc:  svc,
			name: elbv2.LoadBalancerName,
			arn:  elbv2.LoadBalancerArn,
		})
	}

	return resources, nil
}

func (e *ELBv2LoadBalancer) Remove() error {
	params := &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: e.arn,
	}

	_, err := e.svc.DeleteLoadBalancer(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *ELBv2LoadBalancer) String() string {
	return *e.name
}
