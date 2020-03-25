package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/fms"
)

type FMSPolicy struct {
	svc      *fms.FMS
	policyId *string
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
				svc:      svc,
				policyId: policy.PolicyId,
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
		PolicyId:                 f.policyId,
		DeleteAllPolicyResources: aws.Bool(false),
	})

	return err
}

func (f *FMSPolicy) String() string {
	return *f.policyId
}
