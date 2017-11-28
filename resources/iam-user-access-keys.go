package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMUserAccessKey struct {
	svc         *iam.IAM
	accessKeyId string
	userName    string
	status      string
}

func (n *IAMNuke) ListUserAccessKeys() ([]Resource, error) {
	resp, err := n.Service.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Users {
		resp, err := n.Service.ListAccessKeys(
			&iam.ListAccessKeysInput{
				UserName: role.UserName,
			})
		if err != nil {
			return nil, err
		}

		for _, meta := range resp.AccessKeyMetadata {
			resources = append(resources, &IAMUserAccessKey{
				svc:         n.Service,
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
