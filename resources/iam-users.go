package resources

import "github.com/aws/aws-sdk-go/service/iam"

type IAMUser struct {
	svc  *iam.IAM
	name string
}

func (n *IAMNuke) ListUsers() ([]Resource, error) {
	resp, err := n.Service.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Users {
		resources = append(resources, &IAMUser{
			svc:  n.Service,
			name: *out.UserName,
		})
	}

	return resources, nil
}

func (e *IAMUser) Remove() error {
	_, err := e.svc.DeleteUser(&iam.DeleteUserInput{
		UserName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUser) String() string {
	return e.name
}
