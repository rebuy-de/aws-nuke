package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamRolePolicyAttachement struct {
	svc        *iam.IAM
	policyArn  string
	policyName string
	roleName   string
}

func (n *IamNuke) ListRolePolicyAttachements() ([]Resource, error) {
	resp, err := n.Service.ListRoles(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Roles {
		resp, err := n.Service.ListAttachedRolePolicies(
			&iam.ListAttachedRolePoliciesInput{
				RoleName: role.RoleName,
			})
		if err != nil {
			return nil, err
		}

		for _, pol := range resp.AttachedPolicies {
			resources = append(resources, &IamRolePolicyAttachement{
				svc:        n.Service,
				policyArn:  *pol.PolicyArn,
				policyName: *pol.PolicyName,
				roleName:   *role.RoleName,
			})
		}
	}

	return resources, nil
}

func (e *IamRolePolicyAttachement) Filter() error {
	if strings.HasPrefix(e.policyArn, "arn:aws:iam::aws:policy/aws-service-role/") {
		return fmt.Errorf("cannot detach from service roles")
	}
	return nil
}

func (e *IamRolePolicyAttachement) Remove() error {
	_, err := e.svc.DetachRolePolicy(
		&iam.DetachRolePolicyInput{
			PolicyArn: &e.policyArn,
			RoleName:  &e.roleName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamRolePolicyAttachement) String() string {
	return fmt.Sprintf("%s -> %s", e.roleName, e.policyName)
}
