package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMGroupPolicy struct {
	svc        *iam.IAM
	policyName string
	groupName  string
}

func init() {
	register("IAMGroupPolicy", ListIAMGroupPolicies)
}

func ListIAMGroupPolicies(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListGroups(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, group := range resp.Groups {
		resp, err := svc.ListGroupPolicies(
			&iam.ListGroupPoliciesInput{
				GroupName: group.GroupName,
			})
		if err != nil {
			return nil, err
		}

		for _, pol := range resp.PolicyNames {
			resources = append(resources, &IAMGroupPolicy{
				svc:        svc,
				policyName: *pol,
				groupName:  *group.GroupName,
			})
		}
	}

	return resources, nil
}

func (e *IAMGroupPolicy) Remove() error {
	_, err := e.svc.DeleteGroupPolicy(
		&iam.DeleteGroupPolicyInput{
			PolicyName: &e.policyName,
			GroupName:  &e.groupName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMGroupPolicy) String() string {
	return fmt.Sprintf("%s -> %s", e.groupName, e.policyName)
}
