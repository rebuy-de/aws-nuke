package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamInstanceProfileRole struct {
	svc     *iam.IAM
	role    string
	profile string
}

func (n *IamNuke) ListInstanceProfileRoles() ([]Resource, error) {
	resp, err := n.Service.ListInstanceProfiles(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.InstanceProfiles {
		for _, role := range out.Roles {
			resources = append(resources, &IamInstanceProfileRole{
				svc:     n.Service,
				profile: *out.InstanceProfileName,
				role:    *role.RoleName,
			})
		}
	}

	return resources, nil
}

func (e *IamInstanceProfileRole) Remove() error {
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

func (e *IamInstanceProfileRole) String() string {
	return fmt.Sprintf("%s -> %s", e.profile, e.role)
}
