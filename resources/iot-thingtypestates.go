package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTThingTypeState struct {
	svc             *iot.IoT
	name            *string
	deprecated      *bool
	deprecatedEpoch *time.Time
}

func init() {
	register("IoTThingTypeState", ListIoTThingTypeStates)
}

func ListIoTThingTypeStates(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &IoTThingTypeState{
				svc:             svc,
				name:            thingType.ThingTypeName,
				deprecated:      thingType.ThingTypeMetadata.Deprecated,
				deprecatedEpoch: thingType.ThingTypeMetadata.DeprecationDate,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTThingTypeState) Remove() error {

	_, err := f.svc.DeprecateThingType(&iot.DeprecateThingTypeInput{
		ThingTypeName: f.name,
	})

	return err
}

func (f *IoTThingTypeState) String() string {
	return *f.name
}

func (f *IoTThingTypeState) Filter() error {
	//Ensure we don't inspect time unless its already deprecated
	if *f.deprecated == true {
		currentTime := time.Now()
		timeDiff := currentTime.Sub(*f.deprecatedEpoch)
		// Must wait for 300 seconds before deleting a ThingType after deprecation
		// Padding 5 seconds to ensure we are beyond any skew
		if timeDiff < 305 {
			return fmt.Errorf("already deprecated")
		}
	}
	return nil
}
