package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMUserPolicy struct {
	svc        *iam.IAM
	userName   string
	policyName string
}

func init() {
	register("IAMUserPolicy", ListIAMUserPolicies)
}

func ListIAMUserPolicies(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	users, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, user := range users.Users {
		policies, err := svc.ListUserPolicies(&iam.ListUserPoliciesInput{
			UserName: user.UserName,
		})
		if err != nil {
			return nil, err
		}

		for _, policyName := range policies.PolicyNames {
			resources = append(resources, &IAMUserPolicy{
				svc:        svc,
				policyName: *policyName,
				userName:   *user.UserName,
			})
		}
	}

	return resources, nil
}

func (e *IAMUserPolicy) Remove() error {
	_, err := e.svc.DeleteUserPolicy(
		&iam.DeleteUserPolicyInput{
			UserName:   &e.userName,
			PolicyName: &e.policyName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUserPolicy) String() string {
	return fmt.Sprintf("%s -> %s", e.userName, e.policyName)
}
