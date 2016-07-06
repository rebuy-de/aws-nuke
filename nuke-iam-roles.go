package main

import "github.com/aws/aws-sdk-go/service/iam"

type IamRole struct {
	svc  *iam.IAM
	name string
}

func (n *IamNuke) ListRoles() ([]Resource, error) {
	resp, err := n.svc.ListRoles(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Roles {
		resources = append(resources, &IamRole{
			svc:  n.svc,
			name: *out.RoleName,
		})
	}

	return resources, nil
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
