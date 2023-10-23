package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/scheduler"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EventBridgeSchedule struct {
	svc          *scheduler.Scheduler
	scheduleName *string
}

func init() {
	register("EventBridgeSchedule", ListEventBridgeSchedules)
}

func ListEventBridgeSchedules(sess *session.Session) ([]Resource, error) {
	svc := scheduler.New(sess)
	resources := []Resource{}

	params := &scheduler.ListSchedulesInput{}

	for {
		resp, err := svc.ListSchedules(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Schedules {
			resources = append(resources, &EventBridgeSchedule{
				svc:          svc,
				scheduleName: item.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

// The ClientToken current must be generated manually because of a bug in the DeleteSchedule call:
// - https://github.com/aws/aws-sdk-go/issues/4701
func (f *EventBridgeSchedule) Remove() error {

	_, err := f.svc.DeleteSchedule(&scheduler.DeleteScheduleInput{
		Name:        f.scheduleName,
		ClientToken: aws.String("aws-nuke"),
	})

	return err
}

func (f *EventBridgeSchedule) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("scheduleName", f.scheduleName)
	return properties
}
