package resources

import "github.com/aws/aws-sdk-go/service/iam"

type IamGroup struct {
	svc  *iam.IAM
	name string
}

func (n *IamNuke) ListGroups() ([]Resource, error) {
	resp, err := n.Service.ListGroups(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Groups {
		resources = append(resources, &IamGroup{
			svc:  n.Service,
			name: *out.GroupName,
		})
	}

	return resources, nil
}

func (e *IamGroup) Remove() error {
	_, err := e.svc.DeleteGroup(&iam.DeleteGroupInput{
		GroupName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamGroup) String() string {
	return e.name
}
