package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type IAMRole struct {
	svc  *iam.IAM
	role *iam.Role
	name string
	path string
}

func init() {
	register("IAMRole", ListIAMRoles)
}

func GetIAMRole(svc *iam.IAM, roleName *string) (*iam.Role, error) {
	params := &iam.GetRoleInput{
		RoleName: roleName,
	}
	resp, err := svc.GetRole(params)
	return resp.Role, err
}

func ListIAMRoles(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	params := &iam.ListRolesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListRoles(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.Roles {
			role, err := GetIAMRole(svc, out.RoleName)
			if err != nil {
				logrus.
					WithError(err).
					WithField("roleName", *out.RoleName).
					Error("Failed to get listed role")
				continue
			}

			resources = append(resources, &IAMRole{
				svc:  svc,
				role: role,
				name: *role.RoleName,
				path: *role.Path,
			})
		}

		if *resp.IsTruncated == false {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (e *IAMRole) Filter() error {
	if strings.HasPrefix(e.path, "/aws-service-role/") {
		return fmt.Errorf("cannot delete service roles")
	}
	return nil
}

func (e *IAMRole) Remove() error {
	_, err := e.svc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (role *IAMRole) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range role.role.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.
		Set("Name", role.name).
		Set("Path", role.path)
	return properties
}

func (e *IAMRole) String() string {
	return e.name
}
