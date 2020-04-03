package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/pkg/types"
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
			getroleParams := &iam.GetRoleInput{
				RoleName: out.RoleName,
			}
			getroleOutput, err := svc.GetRole(getroleParams)
			if err != nil {
				return nil, err
			}
			resources = append(resources, &IAMRole{
				svc:  svc,
				role: getroleOutput.Role,
				name: *getroleOutput.Role.RoleName,
				path: *getroleOutput.Role.Path,
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
	properties.Set("Name", role.name)
	return properties
}

func (e *IAMRole) String() string {
	return e.name
}
