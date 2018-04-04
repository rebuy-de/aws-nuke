package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTThingGroup struct {
	svc     *iot.IoT
	name    *string
	version *int64
}

func init() {
	register("IoTThingGroup", ListIoTThingGroups)
}

func ListIoTThingGroups(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}
	thingGroups := []*iot.GroupNameAndArn{}

	params := &iot.ListThingGroupsInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListThingGroups(params)
		if err != nil {
			return nil, err
		}

		for _, thingGroup := range output.ThingGroups {
			thingGroups = append(thingGroups, thingGroup)
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	for _, thingGroup := range thingGroups {
		output, err := svc.DescribeThingGroup(&iot.DescribeThingGroupInput{
			ThingGroupName: thingGroup.GroupName,
		})
		if err != nil {
			return nil, err
		}

		resources = append(resources, &IoTThingGroup{
			svc:     svc,
			name:    thingGroup.GroupName,
			version: output.Version,
		})
	}

	return resources, nil
}

func (f *IoTThingGroup) Remove() error {

	_, err := f.svc.DeleteThingGroup(&iot.DeleteThingGroupInput{
		ThingGroupName:  f.name,
		ExpectedVersion: f.version,
	})

	return err
}

func (f *IoTThingGroup) String() string {
	return *f.name
}
