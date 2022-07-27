package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func init() {
	register("CloudWatchEventsBuses", ListCloudWatchEventsBuses)
}

func ListCloudWatchEventsBuses(sess *session.Session) ([]Resource, error) {
	svc := cloudwatchevents.New(sess)

	resp, err := svc.ListEventBuses(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, bus := range resp.EventBuses {
		if *bus.Name == "default" {
			continue
		}

		resources = append(resources, &CloudWatchEventsBus{
			svc:  svc,
			name: bus.Name,
		})
	}
	return resources, nil
}

type CloudWatchEventsBus struct {
	svc  *cloudwatchevents.CloudWatchEvents
	name *string
}

func (bus *CloudWatchEventsBus) Remove() error {
	_, err := bus.svc.DeleteEventBus(&cloudwatchevents.DeleteEventBusInput{
		Name: bus.name,
	})
	return err
}

func (bus *CloudWatchEventsBus) String() string {
	return *bus.name
}
