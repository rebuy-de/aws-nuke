package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMMaintenanceWindow struct {
	svc *ssm.SSM
	ID  *string
}

func init() {
	register("SSMMaintenanceWindow", ListSSMMaintenanceWindows)
}

func ListSSMMaintenanceWindows(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	params := &ssm.DescribeMaintenanceWindowsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeMaintenanceWindows(params)
		if err != nil {
			return nil, err
		}

		for _, windowIdentity := range output.WindowIdentities {
			resources = append(resources, &SSMMaintenanceWindow{
				svc: svc,
				ID:  windowIdentity.WindowId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMMaintenanceWindow) Remove() error {

	_, err := f.svc.DeleteMaintenanceWindow(&ssm.DeleteMaintenanceWindowInput{
		WindowId: f.ID,
	})

	return err
}

func (f *SSMMaintenanceWindow) String() string {
	return *f.ID
}
