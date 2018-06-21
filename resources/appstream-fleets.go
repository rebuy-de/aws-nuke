package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appstream"
)

type AppStreamFleet struct {
	svc  *appstream.AppStream
	name *string
}

func init() {
	register("AppStreamFleet", ListAppStreamFleets)
}

func ListAppStreamFleets(sess *session.Session) ([]Resource, error) {
	svc := appstream.New(sess)
	resources := []Resource{}

	params := &appstream.DescribeFleetsInput{}

	for {
		output, err := svc.DescribeFleets(params)
		if err != nil {
			return nil, err
		}

		for _, fleet := range output.Fleets {
			resources = append(resources, &AppStreamFleet{
				svc:  svc,
				name: fleet.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AppStreamFleet) Remove() error {

	_, err := f.svc.StopFleet(&appstream.StopFleetInput{
		Name: f.name,
	})

	if err != nil {
		return err
	}

	_, err = f.svc.DeleteFleet(&appstream.DeleteFleetInput{
		Name: f.name,
	})

	return err
}

func (f *AppStreamFleet) String() string {
	return *f.name
}
