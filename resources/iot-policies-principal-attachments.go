package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTPolicyPrincipalAttachments struct {
	svc        *iot.IoT
	policyName *string
	principal  *string
}

func init() {
	register("IoTPolicyPrincipalAttachments", ListIoTPolicyPrincipalAttachments)
}

func ListIoTPolicyPrincipalAttachments(sess *session.Session) ([]Resource, error) {
	resources := []Resource{}

	results, err := ListIoTPolicies(sess)
	if err != nil {
		return nil, err
	}

	for _, resource := range results {
		iotPolicy := resource.(*IoTPolicy)

		targets, err := iotPolicy.ListPolicyTargets()
		if err != nil {
			return nil, err
		}
		resources = append(resources, targets...)
	}

	return resources, nil
}

func (f *IoTPolicy) ListPolicyTargets() ([]Resource, error) {
	resources := make([]Resource, 0)
	params := &iot.ListTargetsForPolicyInput{
		PolicyName: f.name,
		PageSize:   aws.Int64(25),
	}
	for {
		output, err := f.svc.ListTargetsForPolicy(params)
		if err != nil {
			return nil, err
		}

		for _, target := range output.Targets {
			resources = append(resources, &IoTPolicyPrincipalAttachments{
				svc:        f.svc,
				policyName: f.name,
				principal:  target,
			})
		}
		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}
	return resources, nil
}

func (f *IoTPolicyPrincipalAttachments) Remove() error {
	_, err := f.svc.DetachPrincipalPolicy(&iot.DetachPrincipalPolicyInput{
		PolicyName: f.policyName,
		Principal:  f.principal,
	})
	return err
}

func (f *IoTPolicyPrincipalAttachments) String() string {
	return fmt.Sprintf("%s -> %s", *f.policyName, *f.principal)
}
