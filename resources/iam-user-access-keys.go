package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMUserAccessKey struct {
	svc         *iam.IAM
	accessKeyId string
	userName    string
	status      string
}

func init() {
	register("IAMUserAccessKey", ListIAMUserAccessKeys)
}

func ListIAMUserAccessKeys(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Users {
		resp, err := svc.ListAccessKeys(
			&iam.ListAccessKeysInput{
				UserName: role.UserName,
			})
		if err != nil {
			return nil, err
		}

		for _, meta := range resp.AccessKeyMetadata {
			resources = append(resources, &IAMUserAccessKey{
				svc:         svc,
				accessKeyId: *meta.AccessKeyId,
				userName:    *meta.UserName,
				status:      *meta.Status,
			})
		}
	}

	return resources, nil
}

func (e *IAMUserAccessKey) Remove() error {
	_, err := e.svc.DeleteAccessKey(
		&iam.DeleteAccessKeyInput{
			AccessKeyId: &e.accessKeyId,
			UserName:    &e.userName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUserAccessKey) String() string {
	return fmt.Sprintf("%s -> %s", e.userName, e.accessKeyId)
}
