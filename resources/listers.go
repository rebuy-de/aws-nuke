package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetListers(sess *session.Session) []ResourceLister {
	var (
		autoscaling = AutoScalingNuke{autoscaling.New(sess)}
		ec2         = EC2Nuke{ec2.New(sess)}
		elb         = ElbNuke{elb.New(sess)}
		route53     = Route53Nuke{route53.New(sess)}
		s3          = S3Nuke{s3.New(sess)}
		iam         = IamNuke{iam.New(sess)}
	)

	return []ResourceLister{
		elb.ListELBs,

		autoscaling.ListGroups,

		route53.ListResourceRecords,
		route53.ListHostedZones,

		ec2.ListKeyPairs,
		ec2.ListInstances,
		ec2.ListSecurityGroups,
		ec2.ListVpnGatewayAttachements,
		ec2.ListVpnConnections,
		ec2.ListNetworkACLs,
		ec2.ListSubnets,
		ec2.ListCustomerGateways,
		ec2.ListVpnGateways,
		ec2.ListInternetGatewayAttachements,
		ec2.ListInternetGateways,
		ec2.ListRouteTables,
		ec2.ListDhcpOptions,
		ec2.ListVpcs,

		iam.ListInstanceProfileRoles,
		iam.ListInstanceProfiles,
		iam.ListRolePolicyAttachements,
		iam.ListRoles,

		s3.ListObjects,
		s3.ListBuckets,
	}
}
