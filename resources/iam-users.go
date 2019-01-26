package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMUser struct {
	svc  *iam.IAM
	name string
}

func init() {
	register("IAMUser", ListIAMUsers)
}

func ListIAMUsers(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Users {
		resources = append(resources, &IAMUser{
			svc:  svc,
			name: *out.UserName,
		})
	}

	return resources, nil
}

func (e *IAMUser) Remove() error {
	err := e.RemoveAllMFADevices()
	if err != nil {
		return err
	}

	err = e.RemoveAllSSHPublicKeys()
	if err != nil {
		return err
	}

	_, err = e.svc.DeleteUser(&iam.DeleteUserInput{
		UserName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUser) RemoveAllMFADevices() error {
	params := &iam.ListMFADevicesInput{
		UserName: aws.String(e.name),
	}

	var devices []*iam.MFADevice
	err := e.svc.ListMFADevicesPages(params,
		func(page *iam.ListMFADevicesOutput, lastPage bool) bool {
			devices = append(devices, page.MFADevices...)
			return true
		})

	if err != nil {
		return err
	}
	for _, device := range devices {
		_, err := e.svc.DeactivateMFADevice(&iam.DeactivateMFADeviceInput{
			UserName:     device.UserName,
			SerialNumber: device.SerialNumber,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *IAMUser) RemoveAllSSHPublicKeys() error {
	params := &iam.ListSSHPublicKeysInput{
		UserName: aws.String(e.name),
	}
	var keys []*iam.SSHPublicKeyMetadata
	err := e.svc.ListSSHPublicKeysPages(params,
		func(page *iam.ListSSHPublicKeysOutput, lastPage bool) bool {
			keys = append(keys, page.SSHPublicKeys...)
			return true
		})

	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err := e.svc.DeleteSSHPublicKey(&iam.DeleteSSHPublicKeyInput{
			UserName:       key.UserName,
			SSHPublicKeyId: key.SSHPublicKeyId,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *IAMUser) String() string {
	return e.name
}
