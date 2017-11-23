package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamUserGroupAttachment struct {
	svc       *iam.IAM
	groupName string
	userName  string
}

func (n *IamNuke) ListUserGroupAttachments() ([]Resource, error) {
	resp, err := n.Service.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Users {
		resp, err := n.Service.ListGroupsForUser(
			&iam.ListGroupsForUserInput{
				UserName: role.UserName,
			})
		if err != nil {
			return nil, err
		}

		for _, grp := range resp.Groups {
			resources = append(resources, &IamUserGroupAttachment{
				svc:       n.Service,
				groupName: *grp.GroupName,
				userName:  *role.UserName,
			})
		}
	}

	return resources, nil
}

func (e *IamUserGroupAttachment) Remove() error {
	_, err := e.svc.RemoveUserFromGroup(
		&iam.RemoveUserFromGroupInput{
			GroupName: &e.groupName,
			UserName:  &e.userName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamUserGroupAttachment) String() string {
	return fmt.Sprintf("%s -> %s", e.userName, e.groupName)
}
