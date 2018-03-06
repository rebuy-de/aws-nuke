package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type CloudWatchDashboard struct {
	svc           *cloudwatch.CloudWatch
	dashboardName *string
}

func init() {
	register("CloudWatchDashboard", ListCloudWatchDashboards)
}

func ListCloudWatchDashboards(sess *session.Session) ([]Resource, error) {
	svc := cloudwatch.New(sess)
	resources := []Resource{}

	params := &cloudwatch.ListDashboardsInput{}

	for {
		output, err := svc.ListDashboards(params)
		if err != nil {
			return nil, err
		}

		for _, dashboardEntry := range output.DashboardEntries {
			resources = append(resources, &CloudWatchDashboard{
				svc:           svc,
				dashboardName: dashboardEntry.DashboardName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CloudWatchDashboard) Remove() error {

	_, err := f.svc.DeleteDashboards(&cloudwatch.DeleteDashboardsInput{
		DashboardNames: []*string{f.dashboardName},
	})

	return err
}

func (f *CloudWatchDashboard) String() string {
	return *f.dashboardName
}
