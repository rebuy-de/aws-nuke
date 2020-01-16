package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/robomaker"
)

type RoboMakerSimulationApplication struct {
	svc     *robomaker.RoboMaker
	name    *string
	arn     *string
	version *string
}

func init() {
	register("RoboMakerSimulationApplication", ListRoboMakerSimulationApplications)
}

func ListRoboMakerSimulationApplications(sess *session.Session) ([]Resource, error) {
	svc := robomaker.New(sess)
	resources := []Resource{}

	params := &robomaker.ListSimulationApplicationsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListSimulationApplications(params)
		if err != nil {
			return nil, err
		}

		for _, robotSimulationApplication := range resp.SimulationApplicationSummaries {
			resources = append(resources, &RoboMakerSimulationApplication{
				svc:     svc,
				name:    robotSimulationApplication.Name,
				arn:     robotSimulationApplication.Arn,
				version: robotSimulationApplication.Version,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *RoboMakerSimulationApplication) Remove() error {

	request := robomaker.DeleteSimulationApplicationInput{
		Application: f.arn,
	}
	if f.version != nil && *f.version != "$LATEST" {
		request.ApplicationVersion = f.version
	}
	_, err := f.svc.DeleteSimulationApplication(&request)

	return err
}

func (f *RoboMakerSimulationApplication) String() string {
	return *f.name
}
