package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/devicefarm"
)

type DeviceFarmProject struct {
	svc *devicefarm.DeviceFarm
	ARN *string
}

func init() {
	register("DeviceFarmProject", ListDeviceFarmProjects)
}

func ListDeviceFarmProjects(sess *session.Session) ([]Resource, error) {
	svc := devicefarm.New(sess)
	resources := []Resource{}

	params := &devicefarm.ListProjectsInput{}

	for {
		output, err := svc.ListProjects(params)
		if err != nil {
			return nil, err
		}

		for _, project := range output.Projects {
			resources = append(resources, &DeviceFarmProject{
				svc: svc,
				ARN: project.Arn,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *DeviceFarmProject) Remove() error {

	_, err := f.svc.DeleteProject(&devicefarm.DeleteProjectInput{
		Arn: f.ARN,
	})

	return err
}

func (f *DeviceFarmProject) String() string {
	return *f.ARN
}
