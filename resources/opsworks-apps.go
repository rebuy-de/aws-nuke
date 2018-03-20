package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
)

type OpsWorksApp struct {
	svc *opsworks.OpsWorks
	ID  *string
}

func init() {
	register("OpsWorksApp", ListOpsWorksApps)
}

func ListOpsWorksApps(sess *session.Session) ([]Resource, error) {
	svc := opsworks.New(sess)
	resources := []Resource{}

	stackParams := &opsworks.DescribeStacksInput{}

	resp, err := svc.DescribeStacks(stackParams)
	if err != nil {
		return nil, err
	}

	appsParams := &opsworks.DescribeAppsInput{}

	for _, stack := range resp.Stacks {
		appsParams.StackId = stack.StackId
		output, err := svc.DescribeApps(appsParams)
		if err != nil {
			return nil, err
		}

		for _, app := range output.Apps {
			resources = append(resources, &OpsWorksApp{
				svc: svc,
				ID:  app.AppId,
			})
		}

	}
	return resources, nil
}

func (f *OpsWorksApp) Remove() error {

	_, err := f.svc.DeleteApp(&opsworks.DeleteAppInput{
		AppId: f.ID,
	})

	return err
}

func (f *OpsWorksApp) String() string {
	return *f.ID
}
