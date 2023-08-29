package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RedshiftScheduledAction struct {
	svc                 *redshift.Redshift
	scheduledActionName *string
}

func init() {
	register("RedshiftScheduledAction", ListRedshiftScheduledActions)
}

func ListRedshiftScheduledActions(sess *session.Session) ([]Resource, error) {
	svc := redshift.New(sess)
	resources := []Resource{}

	params := &redshift.DescribeScheduledActionsInput{}

	for {
		resp, err := svc.DescribeScheduledActions(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ScheduledActions {
			resources = append(resources, &RedshiftScheduledAction{
				svc:                 svc,
				scheduledActionName: item.ScheduledActionName,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (f *RedshiftScheduledAction) Remove() error {

	_, err := f.svc.DeleteScheduledAction(&redshift.DeleteScheduledActionInput{
		ScheduledActionName: f.scheduledActionName,
	})

	return err
}

func (f *RedshiftScheduledAction) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("scheduledActionName", f.scheduledActionName)
	return properties
}
