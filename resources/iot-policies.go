package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTPolicy struct {
	svc                *iot.IoT
	name               *string
	targets            []*string
	deprecatedVersions []*string
}

func init() {
	register("IoTPolicy", ListIoTPolicies)
}

func listIoTPolicyTargets(f *IoTPolicy) (*IoTPolicy, error) {
	targets := []*string{}

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

	f.targets = targets

	return f, nil
}

func listIoTPolicyDeprecatedVersions(f *IoTPolicy) (*IoTPolicy, error) {
	deprecatedVersions := []*string{}

	params := &iot.ListPolicyVersionsInput{
		PolicyName: f.name,
	}

	output, err := f.svc.ListPolicyVersions(params)
	if err != nil {
		return nil, err
	}

	for _, policyVersion := range output.PolicyVersions {
		if *policyVersion.IsDefaultVersion != true {
			deprecatedVersions = append(deprecatedVersions, policyVersion.VersionId)
		}
	}

	f.deprecatedVersions = deprecatedVersions

	return f, nil
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
			p := &IoTPolicy{
				svc:  svc,
				name: policy.PolicyName,
			}

			p, err = listIoTPolicyTargets(p)
			if err != nil {
				return nil, err
			}

			p, err = listIoTPolicyDeprecatedVersions(p)
			if err != nil {
				return nil, err
			}

			resources = append(resources, p)
		}
		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}

	return resources, nil
}

func (f *IoTPolicy) Remove() error {
	// detach attached targets first
	for _, target := range f.targets {
		f.svc.DetachPolicy(&iot.DetachPolicyInput{
			PolicyName: f.name,
			Target:     target,
		})
	}

	// delete deprecated versions
	for _, version := range f.deprecatedVersions {
		f.svc.DeletePolicyVersion(&iot.DeletePolicyVersionInput{
			PolicyName:      f.name,
			PolicyVersionId: version,
		})
	}

	_, err := f.svc.DeletePolicy(&iot.DeletePolicyInput{
		PolicyName: f.name,
	})

	return err
}

func (f *IoTPolicy) String() string {
	return *f.name
}
