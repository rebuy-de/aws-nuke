package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type SageMakerUserProfile struct {
	svc             *sagemaker.SageMaker
	domainID        *string
	userProfileName *string
}

func init() {
	register("SageMakerUserProfiles", ListSageMakerUserProfiles)
}

func ListSageMakerUserProfiles(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListUserProfilesInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListUserProfiles(params)
		if err != nil {
			return nil, err
		}

		for _, userProfile := range resp.UserProfiles {
			resources = append(resources, &SageMakerUserProfile{
				svc:             svc,
				domainID:        userProfile.DomainId,
				userProfileName: userProfile.UserProfileName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerUserProfile) Remove() error {
	_, err := f.svc.DeleteUserProfile(&sagemaker.DeleteUserProfileInput{
		DomainId:        f.domainID,
		UserProfileName: f.userProfileName,
	})

	return err
}

func (f *SageMakerUserProfile) String() string {
	return *f.userProfileName
}

func (i *SageMakerUserProfile) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("DomainID", i.domainID).
		Set("UserProfileName", i.userProfileName)
	return properties
}
