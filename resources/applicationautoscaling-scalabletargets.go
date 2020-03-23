package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
)

func init() {
	register("ApplicationAutoScalingScalableTarget", ListApplicationAutoScalingScalableTargets)
}

func ListApplicationAutoScalingScalableTargets(s *session.Session) ([]Resource, error) {
	svc := applicationautoscaling.New(s)

	// https://docs.aws.amazon.com/autoscaling/application/APIReference/API_RegisterScalableTarget.html#autoscaling-RegisterScalableTarget-request-ServiceNamespace
	namespaces := []string{
		"ecs",
		"elasticmapreduce",
		"ec2",
		"appstream",
		"dynamodb",
		"rds",
		"sagemaker",
		"custom-resource",
		"comprehend",
		"lambda",
	}

	resources := make([]Resource, 0)

	for _, namespace := range namespaces {
		params := &applicationautoscaling.DescribeScalableTargetsInput{
			ServiceNamespace: &namespace,
		}

		for {
			resp, err := svc.DescribeScalableTargets(params)
			if err != nil {
				return nil, err
			}

			for _, aasst := range resp.ScalableTargets {
				resources = append(resources, &ApplicationAutoScalingScalableTarget{
					svc:               svc,
					resourceId:        aasst.ResourceId,
					scalableDimension: aasst.ScalableDimension,
					serviceNamespace:  aasst.ServiceNamespace,
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

type ApplicationAutoScalingScalableTarget struct {
	svc               *applicationautoscaling.ApplicationAutoScaling
	resourceId        *string
	scalableDimension *string
	serviceNamespace  *string
}

func (aasst *ApplicationAutoScalingScalableTarget) Remove() error {
	params := &applicationautoscaling.DeregisterScalableTargetInput{
		ResourceId:        aasst.resourceId,
		ScalableDimension: aasst.scalableDimension,
		ServiceNamespace:  aasst.serviceNamespace,
	}

	_, err := aasst.svc.DeregisterScalableTarget(params)
	if err != nil {
		return err
	}

	return nil
}

func (aasst *ApplicationAutoScalingScalableTarget) String() string {
	return *aasst.resourceId
}
