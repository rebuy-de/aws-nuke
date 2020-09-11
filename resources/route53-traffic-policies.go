package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

func init() {
	register("Route53TrafficPolicy", ListRoute53TrafficPolicies)
}

func ListRoute53TrafficPolicies(sess *session.Session) ([]Resource, error) {
	svc := route53.New(sess)
	params := &route53.ListTrafficPoliciesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListTrafficPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, trafficPolicy := range resp.TrafficPolicySummaries {
			instances, err := instancesForPolicy(svc, trafficPolicy.Id, trafficPolicy.LatestVersion)

			if err != nil {
				return nil, fmt.Errorf("failed to get instance for policy %s %w", *trafficPolicy.Id, err)
			}

			resources = append(resources, &Route53TrafficPolicy{
				svc:       svc,
				id:        trafficPolicy.Id,
				name:      trafficPolicy.Name,
				instances: instances,
			})
		}

		if aws.BoolValue(resp.IsTruncated) == false {
			break
		}
		params.TrafficPolicyIdMarker = resp.TrafficPolicyIdMarker
	}

	return resources, nil
}

func instancesForPolicy(svc *route53.Route53, policyID *string, version *int64) ([]*route53.TrafficPolicyInstance, error) {
	var instances []*route53.TrafficPolicyInstance
	params := &route53.ListTrafficPolicyInstancesByPolicyInput{
		TrafficPolicyId:      policyID,
		TrafficPolicyVersion: version,
	}

	for {
		resp, err := svc.ListTrafficPolicyInstancesByPolicy(params)

		if err != nil {
			return nil, err
		}

		for _, instance := range resp.TrafficPolicyInstances {
			instances = append(instances, instance)
		}

		if aws.BoolValue(resp.IsTruncated) {
			break
		}

		params.TrafficPolicyInstanceTypeMarker = resp.TrafficPolicyInstanceTypeMarker
	}
	return instances, nil
}

type Route53TrafficPolicy struct {
	svc       *route53.Route53
	id        *string
	name      *string
	instances []*route53.TrafficPolicyInstance
}

func (tp *Route53TrafficPolicy) Remove() error {
	params := &route53.DeleteTrafficPolicyInput{
		Id: tp.id,
	}

	_, err := tp.svc.DeleteTrafficPolicy(params)
	return err
}

func (tp *Route53TrafficPolicy) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", *tp.id).
		Set("NAME", *tp.name)
}
