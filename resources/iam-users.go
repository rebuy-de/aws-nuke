package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type IAMUser struct {
	svc              *iam.IAM
	createDate       *time.Time
	name             string
	passwordLastUsed *time.Time
}

func init() {
	register("IAMUser", ListIAMUsers)
}

func ListIAMUsers(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, user := range resp.Users {
		resources = append(resources, &IAMUser{
			svc:              svc,
			createDate:       user.CreateDate,
			name:             *user.UserName,
			passwordLastUsed: user.PasswordLastUsed,
		})
	}

	return resources, nil
}

func (e *IAMUser) Properties() types.Properties {
	properties := types.NewProperties()

	if e.createDate != nil {
		properties.Set("CreateDate", e.createDate.Format(time.RFC3339))
	}
	if e.passwordLastUsed != nil {
		properties.Set("PasswordLastUsed", e.passwordLastUsed.Format(time.RFC3339))
	}

	return properties
}

func (e *IAMUser) Remove() error {
	_, err := e.svc.DeleteUser(&iam.DeleteUserInput{
		UserName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUser) String() string {
	return e.name
}
