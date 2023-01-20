## v1.4.0
- Get latest version of rebuy-de/aws-nuke v2.21.2

## v1.3.0
- Merge in latest version of rebuy-de/aws-nuke from master:
  * Add Athena WorkGroups and NameQueries resource
  * Add EKS Fargate profile resource
  * Add EKS nodegroup resource
  * Add Route53 health check resource
  * Add CloudFormation StackSet resource
  * Add support for disabling CloudFormation termination protection
  * Add Firewall Manager resource
  * Add Comprehend resource
  * Add CloudFormation type resource
  * Add IAM user SSH public key support
  * Add image builder resource
  * Add global accelerator resource
  * Add API Gateway v2 API resource
  * Add API Gateway v2 VPC links resource
  * Add WAF regional rulegroup resource
  * Add CloudWatchLogResourcePolicy resource
  * Add WAFv2 resource
  * Add RDS Event Subscription resource
  * Add AWS Lex resource
  * Add additional support Route53
  * Add Transfer Server & User Resources
  * Add CodeStar Notification Rules resource
- Add Egress Only Internet Gateway resource
- Add Kinesis Analytics Application V2 resource
- Add WAF global rule predicates resource
- Add EMR disable termination protection flag
- Upgrade to go 1.15
- Upgrade to aws sdk version v1.34.12

## v1.2.0

- Merge in latest version of rebuy-de/aws-nuke from master:
  * Add support for Athena WorkGroups and NamedQueries (#464)
  * Upgrade to aws sdk version 1.28.12 (#466)
  * Upgrade to aws sdk version 1.28.12

## v1.1.0

- Add support for Athena WorkGroups and NamedQueries

## v1.0.2

Rename go module back to github.com/rebuy-de/aws-nuke, in order to facilitate merges back to the rebuy-de org.

## v1.0.1

Fix for checksum mismatch on sum.golang.org

## v1.0.0

Initial release, forked from https://github.com/rebuy-de/aws-nuke v2.14.0.