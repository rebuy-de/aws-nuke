package main

import (
	"fmt"

	"github.com/rebuy-de/aws-nuke/vendor/github.com/aws/aws-sdk-go/service/ec2"
)

func nukeEC2Instances(ec2svc *ec2.EC2) {
	dii := &ec2.DescribeInstancesInput{}
	dio, err := ec2svc.DescribeInstances(dii)
	assertNoError(err)

	for _, reservation := range dio.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Printf("ec2.Instance %s", *instance.InstanceId)

			if *instance.State.Name != "running" {
				fmt.Println(" ... not running")
			} else {
				delParams := &ec2.TerminateInstancesInput{
					InstanceIds: []*string{
						instance.InstanceId,
					},
				}
				_, err = ec2svc.TerminateInstances(delParams)
				assertNoError(err)
				fmt.Println(" ... delete requested")
			}
		}
	}
}
