package main

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type EC2Nuke struct {
	svc *ec2.EC2
}

type Route53Nuke struct {
	svc *route53.Route53
}

type AutoScalingNuke struct {
	svc *autoscaling.AutoScaling
}

type ElbNuke struct {
	svc *elb.ELB
}

type S3Nuke struct {
	svc *s3.S3
}
