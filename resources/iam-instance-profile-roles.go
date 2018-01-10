package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMInstanceProfileRole struct {
	svc     *iam.IAM
	role    string
	profile string
}

func init() {
	register("IAMInstanceProfileRole", ListIAMInstanceProfileRoles)
}

func ListIAMInstanceProfileRoles(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListInstanceProfiles(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.InstanceProfiles {
		for _, role := range out.Roles {
			resources = append(resources, &IAMInstanceProfileRole{
				svc:     svc,
				profile: *out.InstanceProfileName,
				role:    *role.RoleName,
			})
		}
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

func (e *IAMInstanceProfileRole) String() string {
	return fmt.Sprintf("%s -> %s", e.profile, e.role)
}
