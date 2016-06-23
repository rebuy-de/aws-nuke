package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/rebuy-de/aws-nuke/vendor/github.com/aws/aws-sdk-go/aws"
)

func nukeAutoScalingGroups(svc *autoscaling.AutoScaling) {
	describeParams := &autoscaling.DescribeAutoScalingGroupsInput{}
	describeResp, err := svc.DescribeAutoScalingGroups(describeParams)
	assertNoError(err)

	for _, asg := range describeResp.AutoScalingGroups {
		fmt.Printf("autoscaling.AutoScalingGroup %s", *asg.AutoScalingGroupName)

		delParams := &autoscaling.DeleteAutoScalingGroupInput{
			AutoScalingGroupName: asg.AutoScalingGroupName,
			ForceDelete:          aws.Bool(true),
		}

		_, err := svc.DeleteAutoScalingGroup(delParams)
		assertNoError(err)
		log.Println(" ... delete requested")
	}
}
