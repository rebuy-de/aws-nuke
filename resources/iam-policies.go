package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IamPolicy struct {
	svc *iam.IAM
	arn string
}

func (n *IamNuke) ListPolicies() ([]Resource, error) {
	resp, err := n.Service.ListPolicies(&iam.ListPoliciesInput{
		Scope: aws.String("Local"),
	})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Policies {
		resources = append(resources, &IamPolicy{
			svc: n.Service,
			arn: *out.Arn,
		})
	}

	return resources, nil
}

func (e *IamPolicy) Remove() error {
	_, err := e.svc.DeletePolicy(&iam.DeletePolicyInput{
		PolicyArn: &e.arn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamPolicy) String() string {
	return e.arn
}
