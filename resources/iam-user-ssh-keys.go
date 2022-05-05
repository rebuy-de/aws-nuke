package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type UserSSHKey struct {
	svc      *iam.IAM
	userName string
	sshKeyID string
}

func init() {
	register("IAMUserSSHPublicKey", ListIAMUserSSHPublicKeys)
}

func ListIAMUserSSHPublicKeys(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	usersOutput, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	var resources []Resource
	for _, user := range usersOutput.Users {
		listOutput, err := svc.ListSSHPublicKeys(&iam.ListSSHPublicKeysInput{
			UserName: user.UserName,
		})

		if err != nil {
			return nil, err
		}

		for _, publicKey := range listOutput.SSHPublicKeys {
			resources = append(resources, &UserSSHKey{
				svc:      svc,
				userName: *user.UserName,
				sshKeyID: *publicKey.SSHPublicKeyId,
			})
		}
	}

	return resources, nil
}

func (u *UserSSHKey) Properties() types.Properties {
	return types.NewProperties().
		Set("UserName", u.userName).
		Set("SSHKeyID", u.sshKeyID)
}

func (u *UserSSHKey) String() string {
	return fmt.Sprintf("%s -> %s", u.userName, u.sshKeyID)
}

func (u *UserSSHKey) Remove() error {
	_, err := u.svc.DeleteSSHPublicKey(&iam.DeleteSSHPublicKeyInput{
		UserName:       &u.userName,
		SSHPublicKeyId: &u.sshKeyID,
	})

	return err
}
