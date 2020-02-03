package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/Optum/aws-nuke/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMRolePolicyAttachment struct {
	svc        *iam.IAM
	policyArn  string
	policyName string
	roleName   string
}

func init() {
	register("IAMRolePolicyAttachment", ListIAMRolePolicyAttachments)
}

func ListIAMRolePolicyAttachments(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	roleParams := &iam.ListRolesInput{}
	resources := make([]Resource, 0)

	for {
		roleResp, err := svc.ListRoles(roleParams)
		if err != nil {
			return nil, err
		}

		for _, role := range roleResp.Roles {
			polParams := &iam.ListAttachedRolePoliciesInput{
				RoleName: role.RoleName,
			}

			for {
				polResp, err := svc.ListAttachedRolePolicies(polParams)
				if err != nil {
					logrus.Errorf("failed to list attached policies for role %s: %v",
						*role.RoleName, err)
					break
				}
				for _, pol := range polResp.AttachedPolicies {
					resources = append(resources, &IAMRolePolicyAttachment{
						svc:        svc,
						policyArn:  *pol.PolicyArn,
						policyName: *pol.PolicyName,
						roleName:   *role.RoleName,
					})
				}

				if *polResp.IsTruncated == false {
					break
				}

				polParams.Marker = polResp.Marker
			}
		}

		if *roleResp.IsTruncated == false {
			break
		}

		roleParams.Marker = roleResp.Marker
	}

	return resources, nil
}

func (e *IAMRolePolicyAttachment) Filter() error {
	if strings.HasPrefix(e.policyArn, "arn:aws:iam::aws:policy/aws-service-role/") {
		return fmt.Errorf("cannot detach from service roles")
	}
	return nil
}

func (e *IAMRolePolicyAttachment) Remove() error {
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

func (e *IAMRolePolicyAttachment) Properties() types.Properties {
	return types.NewProperties().
		Set("RoleName", e.roleName).
		Set("PolicyName", e.policyName)
}

func (e *IAMRolePolicyAttachment) String() string {
	return fmt.Sprintf("%s -> %s", e.roleName, e.policyName)
}
