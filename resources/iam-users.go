package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMUser struct {
	svc  *iam.IAM
	name string
	tags []*iam.Tag
}

func init() {
	register("IAMUser", ListIAMUsers)
}

func GetIAMUser(svc *iam.IAM, userName *string) (*iam.User, error) {
	params := &iam.GetUserInput{
		UserName: userName,
	}
	resp, err := svc.GetUser(params)
	return resp.User, err
}

func ListIAMUsers(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Users {
		user, err := GetIAMUser(svc, out.UserName)
		if err != nil {
			logrus.Errorf("Failed to get user %s: %v", *out.UserName, err)
			continue
		}
		resources = append(resources, &IAMUser{
			svc:  svc,
			name: *out.UserName,
			tags: user.Tags,
		})
	}

	return resources, nil
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

func (e *IAMUser) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", e.name)

	for _, tag := range e.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
