package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type IAMGroupPolicyAttachment struct {
	svc        *iam.IAM
	policyArn  string
	policyName string
	roleName   string
}

func init() {
	register("IAMGroupPolicyAttachment", ListIAMGroupPolicyAttachments)
}

func ListIAMGroupPolicyAttachments(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListGroups(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Groups {
		resp, err := svc.ListAttachedGroupPolicies(
			&iam.ListAttachedGroupPoliciesInput{
				GroupName: role.GroupName,
			})
		if err != nil {
			return nil, err
		}

		for _, pol := range resp.AttachedPolicies {
			resources = append(resources, &IAMGroupPolicyAttachment{
				svc:        svc,
				policyArn:  *pol.PolicyArn,
				policyName: *pol.PolicyName,
				roleName:   *role.GroupName,
			})
		}
	}

	return resources, nil
}

func (e *IAMGroupPolicyAttachment) Remove() error {
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

func (e *IAMGroupPolicyAttachment) Properties() types.Properties {
	return types.NewProperties().
		Set("RoleName", e.roleName).
		Set("PolicyName", e.policyName).
		Set("PolicyArn", e.policyArn)
}

func (e *IAMGroupPolicyAttachment) String() string {
	return fmt.Sprintf("%s -> %s", e.roleName, e.policyName)
}
