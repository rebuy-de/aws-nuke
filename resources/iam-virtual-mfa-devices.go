package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMVirtualMFADevice struct {
	svc          *iam.IAM
	user         *iam.User
	serialNumber string
}

func init() {
	register("IAMVirtualMFADevice", ListIAMVirtualMFADevices)
}

func ListIAMVirtualMFADevices(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListVirtualMFADevices(&iam.ListVirtualMFADevicesInput{})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.VirtualMFADevices {
		resources = append(resources, &IAMVirtualMFADevice{
			svc:          svc,
			user:         out.User,
			serialNumber: *out.SerialNumber,
		})
	}

	return resources, nil
}

func (v *IAMVirtualMFADevice) Filter() error {
	if strings.HasSuffix(v.serialNumber, "/root-account-mfa-device") {
		return fmt.Errorf("Cannot delete root MFA device")
	}
	return nil
}

func (v *IAMVirtualMFADevice) Remove() error {
	if v.user != nil {
		_, err := v.svc.DeactivateMFADevice(&iam.DeactivateMFADeviceInput{
			UserName: v.user.UserName, SerialNumber: &v.serialNumber,
		})
		if err != nil {
			return err
		}
	}

	_, err := v.svc.DeleteVirtualMFADevice(&iam.DeleteVirtualMFADeviceInput{
		SerialNumber: &v.serialNumber,
	})
	return err
}

func (v *IAMVirtualMFADevice) String() string {
	return v.serialNumber
}
