package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMUserPolicyAttachment struct {
	svc        *iam.IAM
	policyArn  string
	policyName string
	userName   string
	userTags   []*iam.Tag
}

func init() {
	register("IAMUserPolicyAttachment", ListIAMUserPolicyAttachments)
}

func ListIAMUserPolicyAttachments(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, user := range resp.Users {
		iamUser, err := GetIAMUser(svc, user.UserName)
		if err != nil {
			logrus.Errorf("Failed to get user %s: %v", *user.UserName, err)
			continue
		}

		resp, err := svc.ListAttachedUserPolicies(
			&iam.ListAttachedUserPoliciesInput{
				UserName: user.UserName,
			})
		if err != nil {
			return nil, err
		}

		for _, pol := range resp.AttachedPolicies {
			resources = append(resources, &IAMUserPolicyAttachment{
				svc:        svc,
				policyArn:  *pol.PolicyArn,
				policyName: *pol.PolicyName,
				userName:   *user.UserName,
				userTags:   iamUser.Tags,
			})
		}
	}

	return resources, nil
}

func (e *IAMUserPolicyAttachment) Remove() error {
	_, err := e.svc.DetachUserPolicy(
		&iam.DetachUserPolicyInput{
			PolicyArn: &e.policyArn,
			UserName:  &e.userName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUserPolicyAttachment) Properties() types.Properties {
	properties := types.NewProperties().
		Set("PolicyArn", e.policyArn).
		Set("PolicyName", e.policyName).
		Set("UserName", e.userName)
	for _, tag := range e.userTags {
		properties.SetTagWithPrefix("user", tag.Key, tag.Value)
	}
	return properties
}

func (e *IAMUserPolicyAttachment) String() string {
	return fmt.Sprintf("%s -> %s", e.userName, e.policyName)
}
