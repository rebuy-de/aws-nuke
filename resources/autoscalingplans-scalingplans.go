package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscalingplans"
)

type AutoScalingPlansScalingPlan struct {
	svc                *autoscalingplans.AutoScalingPlans
	scalingPlanName    *string
	scalingPlanVersion *int64
}

func init() {
	register("AutoScalingPlansScalingPlan", ListAutoScalingPlansScalingPlans)
}

func ListAutoScalingPlansScalingPlans(sess *session.Session) ([]Resource, error) {
	svc := autoscalingplans.New(sess)
	svc.ClientInfo.SigningName = "autoscaling-plans"
	resources := []Resource{}

	params := &autoscalingplans.DescribeScalingPlansInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeScalingPlans(params)
		if err != nil {
			return nil, err
		}

		for _, scalingPlan := range output.ScalingPlans {
			resources = append(resources, &AutoScalingPlansScalingPlan{
				svc:                svc,
				scalingPlanName:    scalingPlan.ScalingPlanName,
				scalingPlanVersion: scalingPlan.ScalingPlanVersion,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AutoScalingPlansScalingPlan) Remove() error {

	_, err := f.svc.DeleteScalingPlan(&autoscalingplans.DeleteScalingPlanInput{
		ScalingPlanName:    f.scalingPlanName,
		ScalingPlanVersion: f.scalingPlanVersion,
	})

	return err
}

func (f *AutoScalingPlansScalingPlan) String() string {
	return *f.scalingPlanName
}
