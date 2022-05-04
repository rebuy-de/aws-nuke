package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2DefaultSecurityGroupRule struct {
	svc      *ec2.EC2
	id       *string
	groupId  *string
	isEgress *bool
}

func init() {
	register("EC2DefaultSecurityGroupRule", ListEC2SecurityGroupRules)
}

func ListEC2SecurityGroupRules(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)

	sgFilters := []*ec2.Filter{
		{
			Name: aws.String("group-name"),
			Values: []*string{
				aws.String("default"),
			},
		},
	}
	groupIds := make([]*string, 0)
	sgParams := &ec2.DescribeSecurityGroupsInput{Filters: sgFilters}
	err := svc.DescribeSecurityGroupsPages(sgParams,
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			for _, group := range page.SecurityGroups {
				groupIds = append(groupIds, group.GroupId)
			}
			return !lastPage
		})
	if err != nil {
		return nil, err
	}

	if len(groupIds) == 0 {
		return resources, nil
	}

	sgRuleFilters := []*ec2.Filter{
		{
			Name:   aws.String("group-id"),
			Values: groupIds,
		},
	}
	sgRuleParams := &ec2.DescribeSecurityGroupRulesInput{Filters: sgRuleFilters}
	err = svc.DescribeSecurityGroupRulesPages(sgRuleParams,
		func(page *ec2.DescribeSecurityGroupRulesOutput, lastPage bool) bool {
			for _, rule := range page.SecurityGroupRules {
				resources = append(resources, &EC2DefaultSecurityGroupRule{
					svc:      svc,
					id:       rule.SecurityGroupRuleId,
					groupId:  rule.GroupId,
					isEgress: rule.IsEgress,
				})
			}
			return !lastPage
		})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *EC2DefaultSecurityGroupRule) Remove() error {
	rules := make([]*string, 1)
	rules[0] = r.id
	if *r.isEgress {
		params := &ec2.RevokeSecurityGroupEgressInput{
			GroupId:              r.groupId,
			SecurityGroupRuleIds: rules,
		}
		_, err := r.svc.RevokeSecurityGroupEgress(params)

		if err != nil {
			return err
		}
	} else {
		params := &ec2.RevokeSecurityGroupIngressInput{
			GroupId:              r.groupId,
			SecurityGroupRuleIds: rules,
		}
		_, err := r.svc.RevokeSecurityGroupIngress(params)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *EC2DefaultSecurityGroupRule) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("SecurityGroupId", r.groupId)
	return properties
}

func (r *EC2DefaultSecurityGroupRule) String() string {
	return *r.id
}
