package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTThing struct {
	svc     *iot.IoT
	name    *string
	version *int64
}

func init() {
	register("IoTThing", ListIoTThings)
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

		for _, thing := range output.Things {
			resources = append(resources, &IoTThing{
				svc:     svc,
				name:    thing.ThingName,
				version: thing.Version,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTThing) Remove() error {
	_, err := f.svc.DeleteThing(&iot.DeleteThingInput{
		ThingName:       f.name,
		ExpectedVersion: f.version,
	})

	return err
}

func (f *IoTThing) String() string {
	return *f.name
}
