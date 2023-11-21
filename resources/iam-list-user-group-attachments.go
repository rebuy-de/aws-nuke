package resources

import (
	"fmt"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
)

type IAMUserGroupAttachment struct {
	svc       *iam.IAM
	groupName string
	userName  string
}

func init() {
	register("IAMUserGroupAttachment", ListIAMUserGroupAttachments)
}

func ListIAMUserGroupAttachments(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	resources := []Resource{}

	err := svc.ListUsersPages(nil, func(page *iam.ListUsersOutput, lastPage bool) bool {
		for _, user := range page.Users {
			err := svc.ListGroupsForUserPages(&iam.ListGroupsForUserInput{
				UserName: user.UserName,
			}, func(groupPage *iam.ListGroupsForUserOutput, lastGroupPage bool) bool {
				for _, group := range groupPage.Groups {
					resources = append(resources, &IAMUserGroupAttachment{
						svc:       svc,
						groupName: *group.GroupName,
						userName:  *user.UserName,
					})
				}
				return !lastGroupPage
			})

			if err != nil {
				logrus.Errorf("failed to list groups for user %s: %v", *user.UserName, err)
				return false
			}
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (e *IAMUserGroupAttachment) Remove() error {
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

func (e *IAMUserGroupAttachment) String() string {
	return fmt.Sprintf("%s -> %s", e.userName, e.groupName)
}

func (e *IAMUserGroupAttachment) Properties() types.Properties {
	return types.NewProperties().
		Set("GroupName", e.groupName).
		Set("UserName", e.userName)
}
