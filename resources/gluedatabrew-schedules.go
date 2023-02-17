package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gluedatabrew"
)

type GlueDataBrewSchedules struct {
	svc  *gluedatabrew.GlueDataBrew
	name *string
}

func init() {
	register("GlueDataBrewSchedules", ListGlueDataBrewSchedules)
}

func ListGlueDataBrewSchedules(sess *session.Session) ([]Resource, error) {
	svc := gluedatabrew.New(sess)
	resources := []Resource{}

	params := &gluedatabrew.ListSchedulesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListSchedules(params)
		if err != nil {
			return nil, err
		}

		for _, schedule := range output.Schedules {
			resources = append(resources, &GlueDataBrewSchedules{
				svc:  svc,
				name: schedule.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDataBrewSchedules) Remove() error {
	_, err := f.svc.DeleteSchedule(&gluedatabrew.DeleteScheduleInput{
		Name: f.name,
	})

	return err
}

func (f *GlueDataBrewSchedules) String() string {
	return *f.name
}
