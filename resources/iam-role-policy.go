package resources

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type IAMRolePolicy struct {
	svc        *iam.IAM
	role       iam.Role
	policyName string
}

func init() {
	register("IAMRolePolicy", ListIAMRolePolicies)
}

func ListIAMRolePolicies(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	roleParams := &iam.ListRolesInput{}
	resources := make([]Resource, 0)

	for {
		roles, err := svc.ListRoles(roleParams)
		if err != nil {
			return nil, err
		}

		for _, listedRole := range roles.Roles {
			role, err := GetIAMRole(svc, listedRole.RoleName)
			if err != nil {
				logrus.Errorf("Failed to get listed role %s: %v", *listedRole.RoleName, err)
				continue
			}

			polParams := &iam.ListRolePoliciesInput{
				RoleName: role.RoleName,
			}

			for {
				policies, err := svc.ListRolePolicies(polParams)
				if err != nil {
					logrus.
						WithError(err).
						WithField("roleName", *role.RoleName).
						Error("Failed to list policies")
					break
				}

				for _, policyName := range policies.PolicyNames {
					resources = append(resources, &IAMRolePolicy{
						svc:        svc,
						role:       *role,
						policyName: *policyName,
					})
				}

				if *policies.IsTruncated == false {
					break
				}

				polParams.Marker = policies.Marker
			}
		}

		if *roles.IsTruncated == false {
			break
		}

		roleParams.Marker = roles.Marker
	}

	return resources, nil
}

func (e *IAMRolePolicy) Filter() error {
	if strings.HasPrefix(aws.StringValue(e.role.Path), "/aws-service-role/") {
		return fmt.Errorf("cannot alter service roles")
	}
	return nil
}

func (e *IAMRolePolicy) Remove() error {
	_, err := e.svc.DeleteRolePolicy(
		&iam.DeleteRolePolicyInput{
			RoleName:   e.role.RoleName,
			PolicyName: &e.policyName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMRolePolicy) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("PolicyName", e.policyName)
	properties.Set("role:RoleName", e.role.RoleName)
	properties.Set("role:RoleID", e.role.RoleId)
	properties.Set("role:Path", e.role.Path)

	for _, tagValue := range e.role.Tags {
		properties.SetTagWithPrefix("role", tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *IAMRolePolicy) String() string {
	return fmt.Sprintf("%s -> %s", aws.StringValue(e.role.RoleName), e.policyName)
}
