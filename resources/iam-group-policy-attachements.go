package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamGroupPolicyAttachement struct {
	svc        *iam.IAM
	policyArn  string
	policyName string
	roleName   string
}

func (n *IamNuke) ListGroupPolicyAttachements() ([]Resource, error) {
	resp, err := n.Service.ListGroups(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Groups {
		resp, err := n.Service.ListAttachedGroupPolicies(
			&iam.ListAttachedGroupPoliciesInput{
				GroupName: role.GroupName,
			})
		if err != nil {
			return nil, err
		}

		for _, pol := range resp.AttachedPolicies {
			resources = append(resources, &IamGroupPolicyAttachement{
				svc:        n.Service,
				policyArn:  *pol.PolicyArn,
				policyName: *pol.PolicyName,
				roleName:   *role.GroupName,
			})
		}
	}

	return resources, nil
}

func (e *IamGroupPolicyAttachement) Remove() error {
	_, err := e.svc.DetachGroupPolicy(
		&iam.DetachGroupPolicyInput{
			PolicyArn: &e.policyArn,
			GroupName: &e.roleName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamGroupPolicyAttachement) String() string {
	return fmt.Sprintf("%s -> %s", e.roleName, e.policyName)
}
