package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ELBv2LoadBalancer struct {
	svc          *elbv2.ELBV2
	name         *string
	arn          *string
	tags         []*elbv2.Tag
	featureFlags config.FeatureFlags
}

func init() {
	register("ELBv2", ListELBv2LoadBalancers)
}

func ListELBv2LoadBalancers(sess *session.Session) ([]Resource, error) {
	svc := elbv2.New(sess)
	var tagReqELBv2ARNs []*string
	ELBv2ArnToName := make(map[string]*string)

	err := svc.DescribeLoadBalancersPages(nil,
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, elbv2 := range page.LoadBalancers {
				tagReqELBv2ARNs = append(tagReqELBv2ARNs, elbv2.LoadBalancerArn)
				ELBv2ArnToName[*elbv2.LoadBalancerArn] = elbv2.LoadBalancerName
			}
			return !lastPage
		})

	if err != nil {
		return nil, err
	}

	// Tags for ELBv2s need to be fetched separately
	// We can only specify up to 20 in a single call
	// See: https://github.com/aws/aws-sdk-go/blob/0e8c61841163762f870f6976775800ded4a789b0/service/elbv2/api.go#L5398
	resources := make([]Resource, 0)
	for len(tagReqELBv2ARNs) > 0 {
		requestElements := len(tagReqELBv2ARNs)
		if requestElements > 20 {
			requestElements = 20
		}

		tagResp, err := svc.DescribeTags(&elbv2.DescribeTagsInput{
			ResourceArns: tagReqELBv2ARNs[:requestElements],
		})
		if err != nil {
			return nil, err
		}
		for _, elbv2TagInfo := range tagResp.TagDescriptions {
			resources = append(resources, &ELBv2LoadBalancer{
				svc:  svc,
				name: ELBv2ArnToName[*elbv2TagInfo.ResourceArn],
				arn:  elbv2TagInfo.ResourceArn,
				tags: elbv2TagInfo.Tags,
			})
		}

		// Remove the elements that were queried
		tagReqELBv2ARNs = tagReqELBv2ARNs[requestElements:]
	}
	return resources, nil
}

func (e *ELBv2LoadBalancer) FeatureFlags(ff config.FeatureFlags) {
	e.featureFlags = ff
}

func (e *ELBv2LoadBalancer) Remove() error {
	params := &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: e.arn,
	}

	_, err := e.svc.DeleteLoadBalancer(params)
	if err != nil {
		if e.featureFlags.DisableDeletionProtection.ELBv2 {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "OperationNotPermitted" &&
				awsErr.Message() == "Load balancer '"+*e.arn+"' cannot be deleted because deletion protection is enabled" {
				err = e.DisableProtection()
				if err != nil {
					return err
				}
				_, err := e.svc.DeleteLoadBalancer(params)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return err
	}
	return nil
}

func (e *ELBv2LoadBalancer) DisableProtection() error {
	params := &elbv2.ModifyLoadBalancerAttributesInput{
		LoadBalancerArn: e.arn,
		Attributes: []*elbv2.LoadBalancerAttribute{
			{
				Key:   aws.String("deletion_protection.enabled"),
				Value: aws.String("false"),
			},
		},
	}
	_, err := e.svc.ModifyLoadBalancerAttributes(params)
	if err != nil {
		return err
	}
	return nil
}

func (e *ELBv2LoadBalancer) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *ELBv2LoadBalancer) String() string {
	return *e.name
}
