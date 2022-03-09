package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type IAMUser struct {
	svc  *iam.IAM
	user *iam.User
}

func init() {
	register("IAMUser", ListIAMUsers)
}

func ListIAMUsers(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	params := &iam.ListUsersInput{}
	resp, err := svc.ListUsers(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, smalluser := range resp.Users {
		getuseroutput, err := svc.GetUser(&iam.GetUserInput{
			UserName: smalluser.UserName,
		})
		if err != nil {
			return nil, err
		}

		fullUser := getuseroutput.User

		thing := &IAMUser{
			svc:  svc,
			user: fullUser,
		}

		resources = append(resources, thing) 
	}

	return resources, nil
}

func (e *IAMUser) Remove() error {
	_, err := e.svc.DeleteUser(&iam.DeleteUserInput{
		UserName: e.user.UserName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUser) String() string {
	return *e.user.UserName
}

func (e *IAMUser) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("UserName", e.user.UserName)
	properties.Set("UserId", e.user.UserId)
	properties.Set("Path", e.user.Path)
	properties.Set("Arn", e.user.Arn)
	properties.Set("CreateDate", e.user.CreateDate)
	properties.Set("PasswordLastUsed", e.user.PasswordLastUsed)


	if (e.user.PermissionsBoundary != nil) {
		properties.Set("PermissionsBoundaryArn", e.user.PermissionsBoundary.PermissionsBoundaryArn)
		properties.Set("PermissionsBoundaryType", e.user.PermissionsBoundary.PermissionsBoundaryType)
	}
	

	for _, tag := range e.user.Tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
