package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMRole struct {
	svc  *iam.IAM
	name string
	path string
}

func init() {
	register("IAMRole", ListIAMRoles)
}

func ListIAMRoles(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListRoles(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Roles {
		resources = append(resources, &IAMRole{
			svc:  svc,
			name: *out.RoleName,
			path: *out.Path,
		})
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

func (e *IAMRole) String() string {
	return e.name
}
