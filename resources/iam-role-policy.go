package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMRolePolicy struct {
	svc        *iam.IAM
	roleName   string
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

		for _, role := range roles.Roles {
			polParams := &iam.ListRolePoliciesInput{
				RoleName: role.RoleName,
			}

			for {
				policies, err := svc.ListRolePolicies(polParams)
				if err != nil {
					return nil, err
				}

				for _, policyName := range policies.PolicyNames {
					resources = append(resources, &IAMRolePolicy{
						svc:        svc,
						policyName: *policyName,
						roleName:   *role.RoleName,
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

func (e *IAMRolePolicy) Remove() error {
	_, err := e.svc.DeleteRolePolicy(
		&iam.DeleteRolePolicyInput{
			RoleName:   &e.roleName,
			PolicyName: &e.policyName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMRolePolicy) String() string {
	return fmt.Sprintf("%s -> %s", e.roleName, e.policyName)
}
