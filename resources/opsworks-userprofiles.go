package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
)

type OpsWorksUserProfile struct {
	svc *opsworks.OpsWorks
	ARN *string
}

func init() {
	register("OpsWorksUserProfile", ListOpsWorksUserProfiles)
}

func ListOpsWorksUserProfiles(sess *session.Session) ([]Resource, error) {
	svc := opsworks.New(sess)
	resources := []Resource{}

	params := &opsworks.DescribeUserProfilesInput{}

	output, err := svc.DescribeUserProfiles(params)
	if err != nil {
		return nil, err
	}

	for _, userProfile := range output.UserProfiles {
		resources = append(resources, &OpsWorksUserProfile{
			svc: svc,
			ARN: userProfile.IamUserArn,
		})
	}

	return resources, nil
}

func (f *OpsWorksUserProfile) Remove() error {

	_, err := f.svc.DeleteUserProfile(&opsworks.DeleteUserProfileInput{
		IamUserArn: f.ARN,
	})

	return err
}

func (f *OpsWorksUserProfile) String() string {
	return *f.ARN
}
