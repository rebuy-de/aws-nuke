package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTRoleAlias struct {
	svc       *iot.IoT
	roleAlias *string
}

func init() {
	register("IoTRoleAlias", ListIoTRoleAliases)
}

func ListIoTRoleAliases(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListRoleAliasesInput{
		PageSize: aws.Int64(25),
	}
	for {
		output, err := svc.ListRoleAliases(params)
		if err != nil {
			return nil, err
		}

		for _, roleAlias := range output.RoleAliases {
			resources = append(resources, &IoTRoleAlias{
				svc:       svc,
				roleAlias: roleAlias,
			})
		}
		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}

	return resources, nil
}

func (f *IoTRoleAlias) Remove() error {

	_, err := f.svc.DeleteRoleAlias(&iot.DeleteRoleAliasInput{
		RoleAlias: f.roleAlias,
	})

	return err
}

func (f *IoTRoleAlias) String() string {
	return *f.roleAlias
}
