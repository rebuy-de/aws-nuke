package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamRole struct {
	svc  *iam.IAM
	name string
	path string
}

func (n *IamNuke) ListRoles() ([]Resource, error) {
	resp, err := n.Service.ListRoles(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Roles {
		resources = append(resources, &IamRole{
			svc:  n.Service,
			name: *out.RoleName,
			path: *out.Path,
		})
	}

	return resources, nil
}

func (e *IamRole) Filter() error {
	if strings.HasPrefix(e.path, "/aws-service-role/") {
		return fmt.Errorf("cannot delete service roles")
	}
	return nil
}

func (e *IamRole) Remove() error {
	_, err := e.svc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamRole) String() string {
	return e.name
}
