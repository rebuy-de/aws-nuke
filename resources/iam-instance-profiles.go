package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMInstanceProfile struct {
	svc     *iam.IAM
	profile *iam.InstanceProfile
	name    string
	path    string
}

func init() {
	register("IAMInstanceProfile", ListIAMInstanceProfiles)
}

func GetIAMInstanceProfile(svc *iam.IAM, instanceProfileName *string) (*iam.InstanceProfile, error) {
	params := &iam.GetInstanceProfileInput{
		InstanceProfileName: instanceProfileName,
	}
	resp, err := svc.GetInstanceProfile(params)
	return resp.InstanceProfile, err
}

func ListIAMInstanceProfiles(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	params := &iam.ListInstanceProfilesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListInstanceProfiles(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.InstanceProfiles {
			profile, err := GetIAMInstanceProfile(svc, out.InstanceProfileName)
			if err != nil {
				logrus.
					WithError(err).
					WithField("instanceProfileName", *out.InstanceProfileName).
					Error("Failed to get listed instance profile")
				continue
			}

			resources = append(resources, &IAMInstanceProfile{
				svc:     svc,
				name:    *out.InstanceProfileName,
				path:    *profile.Path,
				profile: profile,
			})
		}

		if *resp.IsTruncated == false {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (e *IAMInstanceProfile) Remove() error {
	_, err := e.svc.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMInstanceProfile) String() string {
	return e.name
}

func (e *IAMInstanceProfile) Properties() types.Properties {
	properties := types.NewProperties()

	for _, tagValue := range e.profile.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	properties.
		Set("Name", e.name).
		Set("Path", e.path)

	return properties
}
