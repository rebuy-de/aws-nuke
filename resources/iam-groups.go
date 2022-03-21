package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type IAMGroup struct {
	svc  *iam.IAM
	id   string
	name string
	path string
}

func init() {
	register("IAMGroup", ListIAMGroups)
}

func ListIAMGroups(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListGroups(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Groups {
		resources = append(resources, &IAMGroup{
			svc:  svc,
			id:   *out.GroupId,
			name: *out.GroupName,
			path: *out.Path,
		})
	}

	return resources, nil
}

func (e *IAMGroup) Remove() error {
	_, err := e.svc.DeleteGroup(&iam.DeleteGroupInput{
		GroupName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMGroup) String() string {
	return e.name
}

func (e *IAMGroup) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", e.name).
		Set("Path", e.path).
		Set("ID", e.id)
}
