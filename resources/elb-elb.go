package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ELBLoadBalancer struct {
	svc  *elb.ELB
	elb  *elb.LoadBalancerDescription
	tags []*elb.Tag
}

func init() {
	register("ELB", ListELBLoadBalancers)
}

func ListELBLoadBalancers(sess *session.Session) ([]Resource, error) {
	resources := make([]Resource, 0)
	elbNames := make([]*string, 0)
	elbNameToRsc := make(map[string]*elb.LoadBalancerDescription)
	svc := elb.New(sess)

	err := svc.DescribeLoadBalancersPages(nil,
		func(page *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, desc := range page.LoadBalancerDescriptions {
				elbNames = append(elbNames, desc.LoadBalancerName)
				elbNameToRsc[*desc.LoadBalancerName] = desc
			}
			return !lastPage
		})

	if err != nil {
		return nil, err
	}

	for len(elbNames) > 0 {
		requestElements := len(elbNames)
		if requestElements > 20 {
			requestElements = 20
		}

		tagResp, err := svc.DescribeTags(&elb.DescribeTagsInput{
			LoadBalancerNames: elbNames[:requestElements],
		})
		if err != nil {
			return nil, err
		}
		for _, elbTagInfo := range tagResp.TagDescriptions {
			elb := elbNameToRsc[*elbTagInfo.LoadBalancerName]
			resources = append(resources, &ELBLoadBalancer{
				svc:  svc,
				elb:  elb,
				tags: elbTagInfo.Tags,
			})
		}

		// Remove the elements that were queried
		elbNames = elbNames[requestElements:]
	}

	return resources, nil
}

func (e *ELBLoadBalancer) Remove() error {
	params := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: e.elb.LoadBalancerName,
	}

	_, err := e.svc.DeleteLoadBalancer(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *ELBLoadBalancer) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreatedTime", e.elb.CreatedTime.Format(time.RFC3339))

	for _, tagValue := range e.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (e *ELBLoadBalancer) String() string {
	return *e.elb.LoadBalancerName
}
