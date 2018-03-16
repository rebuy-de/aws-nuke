package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appstream"
)

type AppStreamFleetState struct {
	svc   *appstream.AppStream
	name  *string
	state *string
}

func init() {
	register("AppStreamFleetState", ListAppStreamFleetStates)
}

func ListAppStreamFleetStates(sess *session.Session) ([]Resource, error) {
	svc := appstream.New(sess)
	resources := []Resource{}

	params := &appstream.DescribeFleetsInput{}

	for {
		output, err := svc.DescribeFleets(params)
		if err != nil {
			return nil, err
		}

		for _, fleet := range output.Fleets {
			resources = append(resources, &AppStreamFleetState{
				svc:   svc,
				name:  fleet.Name,
				state: fleet.State,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AppStreamFleetState) Remove() error {

	_, err := f.svc.StopFleet(&appstream.StopFleetInput{
		Name: f.name,
	})

	return err
}

func (f *AppStreamFleetState) String() string {
	return *f.name
}

func (f *AppStreamFleetState) Filter() error {
	if *f.state == "STOPPED" {
		return fmt.Errorf("already stopped")
	} else if *f.state == "DELETING" {
		return fmt.Errorf("already being deleted")
	}

	return nil
}
