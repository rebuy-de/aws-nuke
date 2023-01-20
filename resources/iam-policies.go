package resources

import (
	"github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type IAMPolicy struct {
	svc      *iam.IAM
	name     string
	policyId string
	arn      string
	path     string
	tags     []*iam.Tag
}

func init() {
	register("IAMPolicy", ListIAMPolicies)
}

func GetIAMPolicy(svc *iam.IAM, policyArn *string) (*iam.Policy, error) {
	params := &iam.GetPolicyInput{
		PolicyArn: policyArn,
	}
	resp, err := svc.GetPolicy(params)
	return resp.Policy, err
}

func ListIAMPolicies(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	params := &iam.ListPoliciesInput{
		Scope: aws.String("Local"),
	}

	policies := make([]*iam.Policy, 0)

	err := svc.ListPoliciesPages(params,
		func(page *iam.ListPoliciesOutput, lastPage bool) bool {
			for _, listedPolicy := range page.Policies {
				policy, err := GetIAMPolicy(svc, listedPolicy.Arn)
				if err != nil {
					logrus.Errorf("Failed to get listed policy %s: %v", *listedPolicy.PolicyName, err)
					continue
				}
				policies = append(policies, policy)
			}
			return true
		})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, out := range policies {
		resources = append(resources, &IAMPolicy{
			svc:      svc,
			name:     *out.PolicyName,
			path:     *out.Path,
			arn:      *out.Arn,
			policyId: *out.PolicyId,
			tags:     out.Tags,
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

func (policy *IAMPolicy) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("Name", policy.name)
	properties.Set("ARN", policy.arn)
	properties.Set("Path", policy.path)
	properties.Set("PolicyID", policy.policyId)
	for _, tag := range policy.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	return properties
}

func (e *IAMPolicy) String() string {
	return e.arn
}
