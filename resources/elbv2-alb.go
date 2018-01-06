package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type ELBv2 struct {
	svc  *elbv2.ELBV2
	name *string
	arn  *string
}

func init() {
	register("Elbv2ELB", ListElbv2ELBs)
}

func ListElbv2ELBs(sess *session.Session) ([]Resource, error) {
	svc := elbv2.New(sess)

	resp, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, elbv2 := range resp.LoadBalancers {
		resources = append(resources, &ELBv2{
			svc:  svc,
			name: elbv2.LoadBalancerName,
			arn:  elbv2.LoadBalancerArn,
		})
	}

	return resources, nil
}

func (e *ELBv2) Remove() error {
	params := &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: e.arn,
	}

	_, err := e.svc.DeleteLoadBalancer(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *ELBv2) String() string {
	return *e.name
}
