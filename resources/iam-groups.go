package resources

import "github.com/aws/aws-sdk-go/service/iam"

type IAMGroup struct {
	svc  *iam.IAM
	name string
}

func (n *IAMNuke) ListGroups() ([]Resource, error) {
	resp, err := n.Service.ListGroups(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Groups {
		resources = append(resources, &IAMGroup{
			svc:  n.Service,
			name: *out.GroupName,
		})
	}

	return resources, nil
}

func (e *IAMGroup) Remove() error {
	_, err := e.svc.DeleteGroup(&iam.DeleteGroupInput{
		GroupName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMGroup) String() string {
	return e.name
}
