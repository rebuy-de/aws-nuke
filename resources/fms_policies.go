package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/fms"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type FMSPolicy struct {
	svc    *fms.FMS
	policy *fms.PolicySummary
}

func init() {
	register("FMSPolicy", ListFMSPolicies)
}

func ListFMSPolicies(sess *session.Session) ([]Resource, error) {
	svc := fms.New(sess)
	resources := []Resource{}

	params := &fms.ListPoliciesInput{
		MaxResults: aws.Int64(50),
	}

	for {
		resp, err := svc.ListPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, policy := range resp.PolicyList {
			resources = append(resources, &FMSPolicy{
				svc:    svc,
				policy: policy,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *FMSPolicy) Remove() error {

	_, err := f.svc.DeletePolicy(&fms.DeletePolicyInput{
		PolicyId:                 f.policy.PolicyId,
		DeleteAllPolicyResources: aws.Bool(false),
	})

	return err
}

func (f *FMSPolicy) String() string {
	return *f.policy.PolicyId
}

func (f *FMSPolicy) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("PolicyID", f.policy.PolicyId)
	properties.Set("PolicyName", f.policy.PolicyName)
	return properties
}
