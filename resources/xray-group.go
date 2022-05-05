package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type XRayGroup struct {
	svc       *xray.XRay
	groupName *string
	groupARN  *string
}

func init() {
	register("XRayGroup", ListXRayGroups)
}

func ListXRayGroups(sess *session.Session) ([]Resource, error) {
	svc := xray.New(sess)
	resources := []Resource{}

	// Get X-Ray Groups
	var xrayGroups []*xray.GroupSummary
	err := svc.GetGroupsPages(
		&xray.GetGroupsInput{},
		func(page *xray.GetGroupsOutput, lastPage bool) bool {
			for _, group := range page.Groups {
				if *group.GroupName != "Default" { // Ignore the Default group as it cannot be removed
					xrayGroups = append(xrayGroups, group)
				}
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	for _, group := range xrayGroups {
		resources = append(resources, &XRayGroup{
			svc:       svc,
			groupName: group.GroupName,
			groupARN:  group.GroupARN,
		})
	}

	return resources, nil
}

func (f *XRayGroup) Remove() error {
	_, err := f.svc.DeleteGroup(&xray.DeleteGroupInput{
		GroupARN: f.groupARN, // Only allowed to pass GroupARN _or_ GroupName to delete request
	})

	return err
}

func (f *XRayGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("GroupName", f.groupName).
		Set("GroupARN", f.groupARN)

	return properties
}
