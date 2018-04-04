package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTAuthorizer struct {
	svc  *iot.IoT
	name *string
}

func init() {
	register("IoTAuthorizer", ListIoTAuthorizers)
}

func ListIoTAuthorizers(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListAuthorizersInput{}

	output, err := svc.ListAuthorizers(params)
	if err != nil {
		return nil, err
	}

	for _, authorizer := range output.Authorizers {
		resources = append(resources, &IoTAuthorizer{
			svc:  svc,
			name: authorizer.AuthorizerName,
		})
	}

	return resources, nil
}

func (f *IoTAuthorizer) Remove() error {

	_, err := f.svc.DeleteAuthorizer(&iot.DeleteAuthorizerInput{
		AuthorizerName: f.name,
	})

	return err
}

func (f *IoTAuthorizer) String() string {
	return *f.name
}
