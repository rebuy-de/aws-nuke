package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type LightsailLoadBalancer struct {
	svc              *lightsail.Lightsail
	loadBalancerName *string
}

func init() {
	register("LightsailLoadBalancer", ListLightsailLoadBalancers)
}

func ListLightsailLoadBalancers(sess *session.Session) ([]Resource, error) {
	svc := lightsail.New(sess)
	resources := []Resource{}

	params := &lightsail.GetLoadBalancersInput{}

	for {
		output, err := svc.GetLoadBalancers(params)
		if err != nil {
			return nil, err
		}

		for _, loadbalancer := range output.LoadBalancers {
			resources = append(resources, &LightsailLoadBalancer{
				svc:              svc,
				loadBalancerName: loadbalancer.Name,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

func (f *LightsailLoadBalancer) Remove() error {

	_, err := f.svc.DeleteLoadBalancer(&lightsail.DeleteLoadBalancerInput{
		LoadBalancerName: f.loadBalancerName,
	})

	return err
}

func (f *LightsailLoadBalancer) String() string {
	return *f.loadBalancerName
}
