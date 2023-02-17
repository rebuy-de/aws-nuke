package resources

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMRolePolicyAttachment struct {
	svc        *iam.IAM
	policyArn  string
	policyName string
	role       *iam.Role
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

		for _, listedRole := range roleResp.Roles {
			role, err := GetIAMRole(svc, listedRole.RoleName)
			if err != nil {
				logrus.Errorf("Failed to get listed role %s: %v", *listedRole.RoleName, err)
				continue
			}

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
						role:       role,
					})
				}

				if !*polResp.IsTruncated {
					break
				}

				polParams.Marker = polResp.Marker
			}
		}

		if !*roleResp.IsTruncated {
			break
		}

		roleParams.Marker = roleResp.Marker
	}

	return resources, nil
}

func (e *IAMRolePolicyAttachment) Filter() error {
	if strings.Contains(e.policyArn, ":iam::aws:policy/aws-service-role/") {
		return fmt.Errorf("cannot detach from service roles")
	}
	return nil
}

func (e *IAMRolePolicyAttachment) Remove() error {
	_, err := e.svc.DetachRolePolicy(
		&iam.DetachRolePolicyInput{
			PolicyArn: &e.policyArn,
			RoleName:  e.role.RoleName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMRolePolicyAttachment) Properties() types.Properties {
	properties := types.NewProperties().
		Set("RoleName", e.role.RoleName).
		Set("RolePath", e.role.Path).
		Set("RoleLastUsed", getLastUsedDate(e.role, time.RFC3339)).
		Set("RoleCreateDate", e.role.CreateDate.Format(time.RFC3339)).
		Set("PolicyName", e.policyName).
		Set("PolicyArn", e.policyArn)

	for _, tag := range e.role.Tags {
		properties.SetTagWithPrefix("role", tag.Key, tag.Value)
	}
	return properties
}

func (e *IAMRolePolicyAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *e.role.RoleName, e.policyName)
}
