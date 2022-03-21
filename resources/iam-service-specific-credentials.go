package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMServiceSpecificCredential struct {
	svc         *iam.IAM
	name        string
	serviceName string
	id          string
	userName    string
}

func init() {
	register("IAMServiceSpecificCredential", ListServiceSpecificCredentials)
}

func ListServiceSpecificCredentials(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	users, usersErr := ListIAMUsers(sess)
	if usersErr != nil {
		return nil, usersErr
	}

	resources := make([]Resource, 0)
	for _, userResource := range users {
		user, ok := userResource.(*IAMUser)
		if !ok {
			logrus.Errorf("Unable to cast IAMUser.")
			continue
		}
		params := &iam.ListServiceSpecificCredentialsInput{
			UserName: &user.name,
		}
		serviceCredentials, err := svc.ListServiceSpecificCredentials(params)
		if err != nil {
			return nil, err
		}

		for _, credential := range serviceCredentials.ServiceSpecificCredentials {
			resources = append(resources, &IAMServiceSpecificCredential{
				svc:         svc,
				name:        *credential.ServiceUserName,
				serviceName: *credential.ServiceName,
				id:          *credential.ServiceSpecificCredentialId,
				userName:    user.name,
			})
		}
	}

	return resources, nil
}

func (e *IAMServiceSpecificCredential) Remove() error {
	params := &iam.DeleteServiceSpecificCredentialInput{
		ServiceSpecificCredentialId: &e.id,
		UserName:                    &e.userName,
	}
	_, err := e.svc.DeleteServiceSpecificCredential(params)
	if err != nil {
		return err
	}
	return nil
}

func (e *IAMServiceSpecificCredential) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ServiceName", e.serviceName)
	properties.Set("ID", e.id)
	return properties
}

func (e *IAMServiceSpecificCredential) String() string {
	return e.userName + " -> " + e.name
}
