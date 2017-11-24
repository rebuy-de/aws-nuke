package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMPolicy struct {
	svc *iam.IAM
	arn string
}

func (n *IAMNuke) ListPolicies() ([]Resource, error) {
	resp, err := n.Service.ListPolicies(&iam.ListPoliciesInput{
		Scope: aws.String("Local"),
	})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Policies {
		resources = append(resources, &IAMPolicy{
			svc: n.Service,
			arn: *out.Arn,
		})
	}

	return resources, nil
}

func (e *IAMPolicy) Remove() error {
	resp, err := e.svc.ListPolicyVersions(&iam.ListPolicyVersionsInput{
		PolicyArn: &e.arn,
	})
	if err != nil {
		return err
	}
	for _, version := range resp.Versions {
		if !*version.IsDefaultVersion {
			_, err = e.svc.DeletePolicyVersion(&iam.DeletePolicyVersionInput{
				PolicyArn: &e.arn,
				VersionId: version.VersionId,
			})
			if err != nil {
				return err
			}

		}
	}
	_, err = e.svc.DeletePolicy(&iam.DeletePolicyInput{
		PolicyArn: &e.arn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMPolicy) String() string {
	return e.arn
}
