package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTPolicy struct {
	svc  *iot.IoT
	name *string
}

func init() {
	register("IoTPolicy", ListIoTPolicies)
}

func ListIoTPolicies(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListPoliciesInput{
		PageSize: aws.Int64(25),
	}
	for {
		output, err := svc.ListPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, policy := range output.Policies {
			resources = append(resources, &IoTPolicy{
				svc:  svc,
				name: policy.PolicyName,
			})
		}
		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}

	return resources, nil
}

func (f *IoTPolicy) Remove() error {
	err := f.RemoveAllAttachments()
	if err != nil {
		return nil
	}

	err = f.RemoveAllNonDefaultPolicyVersions()
	if err != nil {
		return nil
	}

	_, err = f.svc.DeletePolicy(&iot.DeletePolicyInput{
		PolicyName: f.name,
	})

	return err
}

func (f *IoTPolicy) RemoveAllAttachments() error {
	targets, err := f.ListAllPolicyTargets()
	if err != nil {
		return err
	}
	for _, target := range targets {
		_, err = f.svc.DetachPolicy(&iot.DetachPolicyInput{
			PolicyName: f.name,
			Target:     target,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *IoTPolicy) ListAllPolicyTargets() ([]*string, error) {
	var targets []*string
	params := &iot.ListTargetsForPolicyInput{
		PolicyName: f.name,
		PageSize:   aws.Int64(25),
	}
	for {
		output, err := f.svc.ListTargetsForPolicy(params)
		if err != nil {
			return nil, err
		}
		targets = append(targets, output.Targets...)

		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}
	return targets, nil
}

func (f *IoTPolicy) RemoveAllNonDefaultPolicyVersions() error {
	versions, err := f.ListAllPolicyVersions()
	if err != nil {
		return err
	}

	for _, policyVersion := range versions {
		if aws.BoolValue(policyVersion.IsDefaultVersion) {
			continue
		}
		_, err = f.svc.DeletePolicyVersion(&iot.DeletePolicyVersionInput{
			PolicyName:      f.name,
			PolicyVersionId: policyVersion.VersionId,
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func (f *IoTPolicy) ListAllPolicyVersions() ([]*iot.PolicyVersion, error) {
	output, err := f.svc.ListPolicyVersions(&iot.ListPolicyVersionsInput{
		PolicyName: f.name,
	})
	if err != nil {
		return nil, err
	}
	return output.PolicyVersions, nil

}
func (f *IoTPolicy) String() string {
	return *f.name
}
