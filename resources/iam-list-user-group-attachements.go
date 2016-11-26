package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamUserGroupAttachement struct {
	svc       *iam.IAM
	groupName string
	userName  string
}

func (n *IamNuke) ListUserGroupAttachements() ([]Resource, error) {
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
			resources = append(resources, &IamUserGroupAttachement{
				svc:       n.Service,
				groupName: *grp.GroupName,
				userName:  *role.UserName,
			})
		}
	}

	return resources, nil
}

func (e *IamUserGroupAttachement) Remove() error {
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

func (e *IamUserGroupAttachement) String() string {
	return fmt.Sprintf("%s -> %s", e.userName, e.groupName)
}
