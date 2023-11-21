package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
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
	resources := []Resource{}

	err := svc.ListUsersPages(nil, func(page *iam.ListUsersOutput, lastPage bool) bool {
		for _, out := range page.Users {
			policies, err := svc.ListUserPolicies(&iam.ListUserPoliciesInput{
				UserName: out.UserName,
			})
			if err != nil {
				logrus.Errorf("Failed to list policies for user %s: %v", *out.UserName, err)
				continue
			}

			for _, policyName := range policies.PolicyNames {
				resources = append(resources, &IAMUserPolicy{
					svc:        svc,
					policyName: *policyName,
					userName:   *out.UserName,
				})
			}
		}
		return true
	})

	if err != nil {
		return nil, err
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
