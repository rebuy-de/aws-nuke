package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"github.com/aws/aws-sdk-go/service/sts"
)

type OpsWorksUserProfile struct {
	svc        *opsworks.OpsWorks
	ARN        *string
	callingArn *string
}

func init() {
	register("OpsWorksUserProfile", ListOpsWorksUserProfiles)
}

func ListOpsWorksUserProfiles(sess *session.Session) ([]Resource, error) {
	svc := opsworks.New(sess)
	resources := []Resource{}

	identityOutput, err := sts.New(sess).GetCallerIdentity(nil)
	if err != nil {
		return nil, err
	}

	params := &opsworks.DescribeUserProfilesInput{}

	output, err := svc.DescribeUserProfiles(params)
	if err != nil {
		return nil, err
	}

	for _, userProfile := range output.UserProfiles {
		resources = append(resources, &OpsWorksUserProfile{
			svc:        svc,
			callingArn: identityOutput.Arn,
			ARN:        userProfile.IamUserArn,
		})
	}

	return resources, nil
}

func (f *OpsWorksUserProfile) Filter() error {
	if *f.callingArn == *f.ARN {
		return fmt.Errorf("Cannot delete OpsWorksUserProfile of calling User")
	}
	return nil
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
