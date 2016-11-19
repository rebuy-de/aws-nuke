package resources

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type EC2Nuke struct {
	Service *ec2.EC2
}

type Route53Nuke struct {
	Service *route53.Route53
}

type AutoScalingNuke struct {
	Service *autoscaling.AutoScaling
}

type ElbNuke struct {
	Service *elb.ELB
}

type S3Nuke struct {
	Service *s3.S3
}

type IamNuke struct {
	Service *iam.IAM
}
