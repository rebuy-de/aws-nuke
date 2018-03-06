package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type LightsailStaticIP struct {
	svc          *lightsail.Lightsail
	staticIPName *string
}

func init() {
	register("LightsailStaticIP", ListLightsailStaticIPs)
}

func ListLightsailStaticIPs(sess *session.Session) ([]Resource, error) {
	svc := lightsail.New(sess)
	resources := []Resource{}

	params := &lightsail.GetStaticIpsInput{}

	for {
		output, err := svc.GetStaticIps(params)
		if err != nil {
			return nil, err
		}

		for _, staticIP := range output.StaticIps {
			resources = append(resources, &LightsailStaticIP{
				svc:          svc,
				staticIPName: staticIP.Name,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

func (f *LightsailStaticIP) Remove() error {

	_, err := f.svc.ReleaseStaticIp(&lightsail.ReleaseStaticIpInput{
		StaticIpName: f.staticIPName,
	})

	return err
}

func (f *LightsailStaticIP) String() string {
	return *f.staticIPName
}
