package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTThingPrincipalAttachment struct {
	svc       *iot.IoT
	principal *string
	thingName *string
}

func init() {
	register("IoTThingPrincipalAttachment", ListIoTThingPrincipalAttachments)
}

func ListIoTThingPrincipalAttachments(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	results, err := ListIoTThings(sess)
	if err != nil {
		return nil, err
	}

	for _, resource := range results {
		iotThing := resource.(*IoTThing)
		output, err := svc.ListThingPrincipals(&iot.ListThingPrincipalsInput{
			ThingName: iotThing.name,
		})
		if err != nil {
			return nil, err
		}
		for _, principal := range output.Principals {
			resources = append(resources, &IoTThingPrincipalAttachment{
				svc:       svc,
				principal: principal,
				thingName: iotThing.name,
			})
		}

	}

	return resources, nil
}

func (f *IoTThingPrincipalAttachment) Remove() error {
	_, err := f.svc.DetachThingPrincipal(&iot.DetachThingPrincipalInput{
		Principal: f.principal,
		ThingName: f.thingName,
	})

	return err
}

func (f *IoTThingPrincipalAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.thingName, *f.principal)
}
