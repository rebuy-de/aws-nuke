package main

import "github.com/aws/aws-sdk-go/service/iam"

type IamInstanceProfile struct {
	svc  *iam.IAM
	name string
}

func (n *IamNuke) ListInstanceProfiles() ([]Resource, error) {
	resp, err := n.svc.ListInstanceProfiles(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.InstanceProfiles {
		resources = append(resources, &IamInstanceProfile{
			svc:  n.svc,
			name: *out.InstanceProfileName,
		})
	}

	return resources, nil
}

func (e *IamInstanceProfile) Remove() error {
	_, err := e.svc.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamInstanceProfile) String() string {
	return e.name
}
