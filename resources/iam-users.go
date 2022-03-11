package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/pkg/types"

	"time"
)

type IAMUser struct {
	svc  *iam.IAM

	userName 	*string
	userId 		*string
	path 		*string
	arn 		*string
	createDate 		*time.Time
	passwordLastUsed 	*time.Time
	permissionsBoundaryArn 	*string
	permissionsBoundaryType *string
	tags []*iam.Tag
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

		iamuser := &IAMUser{
			svc:  svc,

			userName: fullUser.UserName,
			userId: fullUser.UserId,
			path: fullUser.Path,
			arn: fullUser.Arn,
			createDate: fullUser.CreateDate,
			passwordLastUsed: fullUser.PasswordLastUsed,
			tags: fullUser.Tags,
		}

		if (fullUser.PermissionsBoundary != nil) {
			iamuser.permissionsBoundaryArn = fullUser.PermissionsBoundary.PermissionsBoundaryArn
			iamuser.permissionsBoundaryType = fullUser.PermissionsBoundary.PermissionsBoundaryType
		}

		resources = append(resources, iamuser) 
	}

	return resources, nil
}

func (e *IAMUser) Remove() error {
	_, err := e.svc.DeleteUser(&iam.DeleteUserInput{
		UserName: e.userName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUser) String() string {
	return *e.userName
}

func (e *IAMUser) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("UserName", e.userName)
	properties.Set("UserId", e.userId)
	properties.Set("Path", e.path)
	properties.Set("Arn", e.arn)
	properties.Set("CreateDate", e.createDate)
	properties.Set("PasswordLastUsed", e.passwordLastUsed)
	properties.Set("PermissionsBoundaryArn", e.permissionsBoundaryArn)
	properties.Set("PermissionsBoundaryType", e.permissionsBoundaryType)

	for _, tag := range e.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
