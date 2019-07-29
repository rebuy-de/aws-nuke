package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTThing struct {
	svc        *iot.IoT
	name       *string
	version    *int64
	principals []*string
}

func init() {
	register("IoTThing", ListIoTThings)
}

func listIoTThingPrincipals(f *IoTThing) (*IoTThing, error) {
	params := &iot.ListThingPrincipalsInput{
		ThingName: f.name,
	}

	output, err := f.svc.ListThingPrincipals(params)
	if err != nil {
		return nil, err
	}

	f.principals = output.Principals
	return f, nil
}

func ListIoTThings(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListThingsInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListThings(params)
		if err != nil {
			return nil, err
		}

		// gather dependent principals
		for _, thing := range output.Things {
			t, err := listIoTThingPrincipals(&IoTThing{
				svc:     svc,
				name:    thing.ThingName,
				version: thing.Version,
			})
			if err != nil {
				return nil, err
			}

			resources = append(resources, t)
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTThing) Remove() error {
	// detach attached principals first
	for _, principal := range f.principals {
		f.svc.DetachThingPrincipal(&iot.DetachThingPrincipalInput{
			Principal: principal,
			ThingName: f.name,
		})
	}

	_, err := f.svc.DeleteThing(&iot.DeleteThingInput{
		ThingName:       f.name,
		ExpectedVersion: f.version,
	})

	return err
}

func (f *IoTThing) String() string {
	return *f.name
}
