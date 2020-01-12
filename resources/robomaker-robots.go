package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/robomaker"
)

type RoboMakerRobot struct {
	svc  *robomaker.RoboMaker
	name *string
	arn  *string
}

func init() {
	register("RoboMakerRobot", ListRoboMakerRobots)
}

func ListRoboMakerRobots(sess *session.Session) ([]Resource, error) {
	svc := robomaker.New(sess)
	resources := []Resource{}

	params := &robomaker.ListRobotsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListRobots(params)
		if err != nil {
			return nil, err
		}

		for _, robot := range resp.Robots {
			resources = append(resources, &RoboMakerRobot{
				svc:  svc,
				name: robot.Name,
				arn:  robot.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *RoboMakerRobot) Remove() error {

	_, err := f.svc.DeleteRobot(&robomaker.DeleteRobotInput{
		Robot: f.arn,
	})

	return err
}

func (f *RoboMakerRobot) String() string {
	return *f.name
}
