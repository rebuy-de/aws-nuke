package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
)

func GetListers(sess *session.Session) []ResourceLister {
	var (
		autoscaling      = AutoScalingNuke{autoscaling.New(sess)}
		cloudformation   = CloudFormationNuke{cloudformation.New(sess)}
		cloudwatchevents = CloudWatchEventsNuke{cloudwatchevents.New(sess)}
		ec2              = EC2Nuke{ec2.New(sess)}
		ecr              = ECRNuke{ecr.New(sess)}
		efs              = EFSNuke{efs.New(sess)}
		elasticache      = ElasticacheNuke{elasticache.New(sess)}
		elb              = ElbNuke{elb.New(sess)}
		iam              = IamNuke{iam.New(sess)}
		rds              = RDSNuke{rds.New(sess)}
		route53          = Route53Nuke{route53.New(sess)}
		s3               = S3Nuke{s3.New(sess)}
		sns              = SNSNuke{sns.New(sess)}
	)

	return []ResourceLister{
		autoscaling.ListGroups,
		cloudformation.ListStacks,
		cloudwatchevents.ListRules,
		cloudwatchevents.ListTargets,
		ec2.ListAddresses,
		ec2.ListCustomerGateways,
		ec2.ListDhcpOptions,
		ec2.ListInstances,
		ec2.ListInternetGatewayAttachements,
		ec2.ListInternetGateways,
		ec2.ListKeyPairs,
		ec2.ListNatGateways,
		ec2.ListNetworkACLs,
		ec2.ListRouteTables,
		ec2.ListSecurityGroups,
		ec2.ListSubnets,
		ec2.ListVolumes,
		ec2.ListVpcs,
		ec2.ListVpnConnections,
		ec2.ListVpnGatewayAttachements,
		ec2.ListVpnGateways,
		ecr.ListRepos,
		efs.ListFileSystems,
		efs.ListMountTargets,
		elasticache.ListCacheClusters,
		elasticache.ListSubnetGroups,
		elb.ListELBs,
		iam.ListGroupPolicyAttachements,
		iam.ListGroups,
		iam.ListInstanceProfileRoles,
		iam.ListInstanceProfiles,
		iam.ListPolicies,
		iam.ListRolePolicyAttachements,
		iam.ListRoles,
		iam.ListServerCertificates,
		iam.ListUserAccessKeys,
		iam.ListUserGroupAttachements,
		iam.ListUserGroupAttachements,
		iam.ListUserPolicyAttachements,
		iam.ListUsers,
		rds.ListInstances,
		rds.ListParameterGroups,
		rds.ListSnapshots,
		rds.ListSubnetGroups,
		route53.ListHostedZones,
		route53.ListResourceRecords,
		s3.ListBuckets,
		s3.ListObjects,
		sns.ListSubscriptions,
		sns.ListTopics,
	}
}
