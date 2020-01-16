package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/robomaker"
)

type RoboMakerRobotApplication struct {
	svc     *robomaker.RoboMaker
	name    *string
	arn     *string
	version *string
}

func init() {
	register("RoboMakerRobotApplication", ListRoboMakerRobotApplications)
}

func ListRoboMakerRobotApplications(sess *session.Session) ([]Resource, error) {
	svc := robomaker.New(sess)
	resources := []Resource{}

	params := &robomaker.ListRobotApplicationsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListRobotApplications(params)
		if err != nil {
			return nil, err
		}

		for _, robotApplication := range resp.RobotApplicationSummaries {
			resources = append(resources, &RoboMakerRobotApplication{
				svc:     svc,
				name:    robotApplication.Name,
				arn:     robotApplication.Arn,
				version: robotApplication.Version,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *RoboMakerRobotApplication) Remove() error {

	request := robomaker.DeleteRobotApplicationInput{
		Application: f.arn,
	}
	if f.version != nil && *f.version != "$LATEST" {
		request.ApplicationVersion = f.version
	}

	_, err := f.svc.DeleteRobotApplication(&request)

	return err
}

func (f *RoboMakerRobotApplication) String() string {
	return *f.name
}
