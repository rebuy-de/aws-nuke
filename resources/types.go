package resources

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AutoScalingNuke struct {
	Service *autoscaling.AutoScaling
}

type CloudFormationNuke struct {
	Service *cloudformation.CloudFormation
}

type EC2Nuke struct {
	Service *ec2.EC2
}

type ECRNuke struct {
	Service *ecr.ECR
}

type EFSNuke struct {
	Service *efs.EFS
}

type ElasticacheNuke struct {
	Service *elasticache.ElastiCache
}

type ElbNuke struct {
	Service *elb.ELB
}

type IamNuke struct {
	Service *iam.IAM
}

type RDSNuke struct {
	Service *rds.RDS
}

type Route53Nuke struct {
	Service *route53.Route53
}

type S3Nuke struct {
	Service *s3.S3
}
