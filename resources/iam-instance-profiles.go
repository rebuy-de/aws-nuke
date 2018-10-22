package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMInstanceProfile struct {
	svc  *iam.IAM
	name string
}

func init() {
	register("IAMInstanceProfile", ListIAMInstanceProfiles)
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
			resources = append(resources, &IAMInstanceProfile{
				svc:  svc,
				name: *out.InstanceProfileName,
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
