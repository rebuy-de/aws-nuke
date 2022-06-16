package resources

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMLoginProfile struct {
	svc  *iam.IAM
	name string
}

func init() {
	register("IAMLoginProfile", ListIAMLoginProfiles)
}

func ListIAMLoginProfiles(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Users {
		lpresp, err := svc.GetLoginProfile(&iam.GetLoginProfileInput{UserName: out.UserName})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeNoSuchEntityException:
					// The user does not have a login profile and we do not
					// need to print an error for that.
					continue
				}
			}

			logrus.Errorf("failed to list login profile for user %s: %v",
				*out.UserName, err)
			continue
		}

		if lpresp.LoginProfile != nil {
			resources = append(resources, &IAMLoginProfile{
				svc:  svc,
				name: *out.UserName,
			})
		}
	}

	return resources, nil
}

func (e *IAMLoginProfile) Remove() error {
	_, err := e.svc.DeleteLoginProfile(&iam.DeleteLoginProfileInput{UserName: &e.name})
	if err != nil {
		return err
	}
	return nil
}

func (e *IAMLoginProfile) Properties() types.Properties {
	return types.NewProperties().
		Set("UserName", e.name)
}

func (e *IAMLoginProfile) String() string {
	return e.name
}
