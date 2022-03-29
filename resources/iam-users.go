package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/pkg/types"
	"github.com/sirupsen/logrus"
  
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

//gets full user atributes with GetUser()
func GetIAMUser(svc *iam.IAM, userName *string) (*iam.User, error) {
	params := &iam.GetUserInput{
		UserName: userName,
	}
	resp, err := svc.GetUser(params)
	return resp.User, err
}

func ListIAMUsers(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	params := &iam.ListUsersInput{}
	resp, err := svc.ListUsers(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Users {
    //ListUsers() does not return all parameters for a users.
		user, err := GetIAMUser(svc, out.UserName)
		if err != nil {
			logrus.Errorf("Failed to get user %s: %v", *out.UserName, err)
			continue
		}
    
		iamuser := &IAMUser{
			svc:  svc,

			userName: user.UserName,
			userId: user.UserId,
			path: user.Path,
			arn: user.Arn,
			createDate: user.CreateDate,
			passwordLastUsed: user.PasswordLastUsed,
			tags: user.Tags,
		}

		if (user.PermissionsBoundary != nil) {
			iamuser.permissionsBoundaryArn = user.PermissionsBoundary.PermissionsBoundaryArn
			iamuser.permissionsBoundaryType = user.PermissionsBoundary.PermissionsBoundaryType
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
