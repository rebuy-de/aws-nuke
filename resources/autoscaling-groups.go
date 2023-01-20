package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("AutoScalingGroup", ListAutoscalingGroups)
}

func ListAutoscalingGroups(s *session.Session) ([]Resource, error) {
	svc := autoscaling.New(s)
	resources := make([]Resource, 0)

	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	err := svc.DescribeAutoScalingGroupsPages(params,
		func(page *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			for _, asg := range page.AutoScalingGroups {
				resources = append(resources, &AutoScalingGroup{
					group: asg,
					svc:   svc,
					tags:  asg.Tags,
				})
			}
			return !lastPage
		})

	if err != nil {
		return nil, err
	}

	return resources, nil
}

type AutoScalingGroup struct {
	svc   *autoscaling.AutoScaling
	group *autoscaling.Group
	tags  []*autoscaling.TagDescription
}

func (asg *AutoScalingGroup) Remove() error {
	params := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: asg.group.AutoScalingGroupName,
		ForceDelete:          aws.Bool(true),
	}

	_, err := asg.svc.DeleteAutoScalingGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (asg *AutoScalingGroup) String() string {
	return *asg.group.AutoScalingGroupName
}

func (asg *AutoScalingGroup) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range asg.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	properties.Set("CreatedTime", asg.group.CreatedTime)
	properties.Set("Name", asg.group.AutoScalingGroupName)

	return properties
}
