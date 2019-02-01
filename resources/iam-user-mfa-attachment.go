package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type IAMUserMfaAttachment struct {
	svc          *iam.IAM
	userName     *string
	serialNumber string
}

func init() {
	register("IAMUserMfaAttachment", ListIAMUserMfaAttachments)
}

func ListIAMUserMfaAttachments(sess *session.Session) ([]Resource, error) {
	results, err := ListIAMUsers(sess)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, resource := range results {
		user := resource.(*IAMUser)
		userMfaDevices, err := user.ListAttachedMFADevices()
		if err != nil {
			return nil, err
		}
		resources = append(resources, userMfaDevices...)
	}

	return resources, nil
}

func (f *IAMUser) ListAttachedMFADevices() ([]Resource, error) {
	resources := make([]Resource, 0)
	params := &iam.ListMFADevicesInput{
		UserName: aws.String(f.name),
	}
	for {
		resp, err := f.svc.ListMFADevices(params)
		if err != nil {
			return nil, err
		}

		for _, device := range resp.MFADevices {
			resources = append(resources, &IAMUserMfaAttachment{
				svc:          f.svc,
				userName:     params.UserName,
				serialNumber: *device.SerialNumber,
			})
		}
		if params.Marker == nil {
			break
		}
		params.Marker = resp.Marker
	}
	return resources, nil
}

func (e *IAMUserMfaAttachment) Remove() error {
	_, err := e.svc.DeactivateMFADevice(&iam.DeactivateMFADeviceInput{
		UserName:     e.userName,
		SerialNumber: &e.serialNumber,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUserMfaAttachment) Properties() types.Properties {
	return types.NewProperties().
		Set("UserName", e.userName).
		Set("SerialNumber", e.serialNumber)
}

func (e *IAMUserMfaAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *e.userName, e.serialNumber)
}
