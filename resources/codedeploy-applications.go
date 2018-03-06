package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
)

type CodeDeployApplication struct {
	svc             *codedeploy.CodeDeploy
	applicationName *string
}

func init() {
	register("CodeDeployApplication", ListCodeDeployApplications)
}

func ListCodeDeployApplications(sess *session.Session) ([]Resource, error) {
	svc := codedeploy.New(sess)
	resources := []Resource{}

	params := &codedeploy.ListApplicationsInput{}

	for {
		resp, err := svc.ListApplications(params)
		if err != nil {
			return nil, err
		}

		for _, application := range resp.Applications {
			resources = append(resources, &CodeDeployApplication{
				svc:             svc,
				applicationName: application,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeDeployApplication) Remove() error {

	_, err := f.svc.DeleteApplication(&codedeploy.DeleteApplicationInput{
		ApplicationName: f.applicationName,
	})

	return err
}

func (f *CodeDeployApplication) String() string {
	return *f.applicationName
}
