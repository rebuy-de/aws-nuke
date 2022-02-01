package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type IAMInstanceProfileRole struct {
	svc     *iam.IAM
	role    string
	profile string
	tags    []*iam.Tag
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
			for _, role := range out.Roles {
				resources = append(resources, &IAMInstanceProfileRole{
					svc:     svc,
					profile: *out.InstanceProfileName,
					role:    *role.RoleName,
					tags:    role.Tags,
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
			InstanceProfileName: &e.profile,
			RoleName:            &e.role,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMInstanceProfileRole) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Role", e.role)
	properties.Set("Profile", e.profile)

	for _, tag := range e.tags {
		properties.SetTagWithPrefix("role", tag.Key, tag.Value)
	}

	return properties
}

func (e *IAMInstanceProfileRole) String() string {
	return fmt.Sprintf("%s -> %s", e.profile, e.role)
}
