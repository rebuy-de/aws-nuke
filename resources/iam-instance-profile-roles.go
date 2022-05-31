package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMInstanceProfileRole struct {
	svc     *iam.IAM
	role    string
	profile *iam.InstanceProfile
}

func init() {
	register("IAMInstanceProfileRole", ListIAMInstanceProfileRoles)
}

func ListIAMInstanceProfileRoles(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	params := &iam.ListInstanceProfilesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListInstanceProfiles(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.InstanceProfiles {
			for _, outRole := range out.Roles {
				profile, err := GetIAMInstanceProfile(svc, out.InstanceProfileName)
				if err != nil {
					logrus.
						WithError(err).
						WithField("instanceProfileName", *out.InstanceProfileName).
						Error("Failed to get listed instance profile")
					continue
				}

				resources = append(resources, &IAMInstanceProfileRole{
					svc:     svc,
					role:    *outRole.RoleName,
					profile: profile,
				})
			}
		}

		if *resp.IsTruncated == false {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (e *IAMInstanceProfileRole) Remove() error {
	_, err := e.svc.RemoveRoleFromInstanceProfile(
		&iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: e.profile.InstanceProfileName,
			RoleName:            &e.role,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMInstanceProfileRole) String() string {
	return fmt.Sprintf("%s -> %s", *e.profile.InstanceProfileName, e.role)
}

func (e *IAMInstanceProfileRole) Properties() types.Properties {
	properties := types.NewProperties()

	for _, tagValue := range e.profile.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	properties.
		Set("InstanceProfile", e.profile.InstanceProfileName).
		Set("InstanceRole", e.role)

	return properties
}
