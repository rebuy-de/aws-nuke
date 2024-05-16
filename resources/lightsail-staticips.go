package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"time"
)

type LightsailStaticIP struct {
	svc          *lightsail.Lightsail
	staticIPName *string
	createdAt    *time.Time
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
				createdAt:    staticIP.CreatedAt,
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

func (f *LightsailStaticIP) Properties() types.Properties {
	return types.NewProperties().
		Set("CreatedAt", f.createdAt.Format(time.RFC3339))
}

func (f *LightsailStaticIP) String() string {
	return *f.staticIPName
}
