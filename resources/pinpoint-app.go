package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pinpoint"
)

type PinpointApp struct {
	svc *pinpoint.Pinpoint
	app string
}

func init() {
	register("PinpointApp", ListPinpointApps)
}

func ListPinpointApps(sess *session.Session) ([]Resource, error) {
	svc := pinpoint.New(sess)

	resp, err := svc.GetApps(&pinpoint.GetAppsInput{})
	if err != nil {
		return nil, err
	}

	apps := make([]Resource, 0)
	for _, appResponse := range resp.ApplicationsResponse.Item {
		apps = append(apps, &PinpointApp{
			svc: svc,
			app: aws.StringValue(appResponse.Id),
		})
	}

	return apps, nil
}

func (p *PinpointApp) Remove() error {
	params := &pinpoint.DeleteAppInput{
		ApplicationId: aws.String(p.app),
	}

	_, err := p.svc.DeleteApp(params)
	if err != nil {
		return err
	}

	return nil
}

func (p *PinpointApp) String() string {
	return p.app
}
