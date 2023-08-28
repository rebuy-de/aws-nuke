package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchrum"
)

type CloudWatchRumApp struct {
	svc            *cloudwatchrum.CloudWatchRUM
	appmonitorname *string
}

func init() {
	register("CloudWatchRUMApp", ListCloudWatchRumApp)
}

func ListCloudWatchRumApp(sess *session.Session) ([]Resource, error) {
	svc := cloudwatchrum.New(sess)
	resources := []Resource{}

	params := &cloudwatchrum.ListAppMonitorsInput{}

	for {
		output, err := svc.ListAppMonitors(params)
		if err != nil {
			return nil, err
		}

		for _, appEntry := range output.AppMonitorSummaries {
			resources = append(resources, &CloudWatchRumApp{
				svc:            svc,
				appmonitorname: appEntry.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CloudWatchRumApp) Remove() error {

	_, err := f.svc.DeleteAppMonitor(&cloudwatchrum.DeleteAppMonitorInput{
		Name: f.appmonitorname,
	})

	return err
}

func (f *CloudWatchRumApp) String() string {
	return *f.appmonitorname
}
