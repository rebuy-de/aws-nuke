package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppAutoScaling struct {
	svc       *applicationautoscaling.ApplicationAutoScaling
	target    *applicationautoscaling.ScalableTarget
	id        string
	roleARN   string
	dimension string
	namespace string
}

func init() {
	register("ApplicationAutoScalingScalableTarget", ListApplicationAutoScalingScalableTargets)
}

func ListApplicationAutoScalingScalableTargets(sess *session.Session) ([]Resource, error) {
	svc := applicationautoscaling.New(sess)

	namespaces := applicationautoscaling.ServiceNamespace_Values()

	params := &applicationautoscaling.DescribeScalableTargetsInput{}
	resources := make([]Resource, 0)
	for _, namespace := range namespaces {
		for {
			params.ServiceNamespace = &namespace
			resp, err := svc.DescribeScalableTargets(params)
			if err != nil {
				return nil, err
			}

			for _, out := range resp.ScalableTargets {
				resources = append(resources, &AppAutoScaling{
					svc:       svc,
					target:    out,
					id:        *out.ResourceId,
					roleARN:   *out.RoleARN,
					dimension: *out.ScalableDimension,
					namespace: *out.ServiceNamespace,
				})
			}

			if resp.NextToken == nil {
				break
			}
			params.NextToken = resp.NextToken
		}
	}

	return resources, nil
}

func (a *AppAutoScaling) Remove() error {
	_, err := a.svc.DeregisterScalableTarget(&applicationautoscaling.DeregisterScalableTargetInput{
		ResourceId:        &a.id,
		ScalableDimension: &a.dimension,
		ServiceNamespace:  &a.namespace,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *AppAutoScaling) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ResourceID", a.id)
	properties.Set("ScalableDimension", a.dimension)
	properties.Set("ServiceNamespace", a.namespace)

	return properties
}

func (a *AppAutoScaling) String() string {
	return a.id + ": " + a.dimension
}
