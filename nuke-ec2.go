package main

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2Nuke struct {
	svc *ec2.EC2
}
