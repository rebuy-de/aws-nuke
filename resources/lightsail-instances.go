package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type LightsailInstance struct {
	svc          *lightsail.Lightsail
	instanceName *string
}

func init() {
	register("LightsailInstance", ListLightsailInstances)
}

func ListLightsailInstances(sess *session.Session) ([]Resource, error) {
	svc := lightsail.New(sess)
	resources := []Resource{}

	params := &lightsail.GetInstancesInput{}

	for {
		output, err := svc.GetInstances(params)
		if err != nil {
			return nil, err
		}

		for _, instance := range output.Instances {
			resources = append(resources, &LightsailInstance{
				svc:          svc,
				instanceName: instance.Name,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

func (f *LightsailInstance) Remove() error {
	_, err := f.svc.DeleteInstance(&lightsail.DeleteInstanceInput{
		InstanceName:      f.instanceName,
		ForceDeleteAddOns: aws.Bool(true),
	})
	return err
}

func (f *LightsailInstance) String() string {
	return *f.instanceName
}
