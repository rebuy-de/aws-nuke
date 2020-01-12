package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/robomaker"
)

type RoboMakerFleet struct {
	svc  *robomaker.RoboMaker
	name *string
	arn  *string
}

func init() {
	register("RoboMakerFleet", ListRoboMakerFleets)
}

func ListRoboMakerFleets(sess *session.Session) ([]Resource, error) {
	svc := robomaker.New(sess)
	resources := []Resource{}

	params := &robomaker.ListFleetsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListFleets(params)
		if err != nil {
			return nil, err
		}

		for _, fleet := range resp.FleetDetails {
			resources = append(resources, &RoboMakerFleet{
				svc:  svc,
				name: fleet.Name,
				arn:  fleet.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *RoboMakerFleet) Remove() error {

	_, err := f.svc.DeleteFleet(&robomaker.DeleteFleetInput{
		Fleet: f.arn,
	})

	return err
}

func (f *RoboMakerFleet) String() string {
	return *f.name
}
