package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTThingType struct {
	svc  *iot.IoT
	name *string
}

func init() {
	register("IoTThingType", ListIoTThingTypes)
}

func ListIoTThingTypes(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListThingTypesInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListThingTypes(params)
		if err != nil {
			return nil, err
		}

		for _, thingType := range output.ThingTypes {
			resources = append(resources, &IoTThingType{
				svc:  svc,
				name: thingType.ThingTypeName,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTThingType) Remove() error {

	_, err := f.svc.DeleteThingType(&iot.DeleteThingTypeInput{
		ThingTypeName: f.name,
	})

	return err
}

func (f *IoTThingType) String() string {
	return *f.name
}
