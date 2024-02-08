package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rolesanywhere"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type Profile struct {
	svc  *rolesanywhere.RolesAnywhere
	ProfileId   string
}

func init() {
	register("IAMRolesAnywhereProfiles", ListProfiles)
}

func ListProfiles(sess *session.Session) ([]Resource, error) {
	svc := rolesanywhere.New(sess)

	params := &rolesanywhere.ListProfilesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListProfiles(params)
		if err != nil {
			return nil, err
		}
		for _, profile := range resp.Profiles {
			resources = append(resources, &Profile{
				svc:      svc,
				ProfileId: *profile.ProfileId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (e *Profile) Remove() error {
	_, err := e.svc.DeleteProfile(&rolesanywhere.DeleteProfileInput{
		ProfileId: &e.ProfileId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Profile) String() string {
	return e.ProfileId
}

func (e *Profile) Properties() types.Properties {
	return types.NewProperties().
		Set("ProfileId", e.ProfileId)
}
