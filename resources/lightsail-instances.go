package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type LightsailInstance struct {
	svc          *lightsail.Lightsail
	instanceName *string
	tags         []*lightsail.Tag

	featureFlags config.FeatureFlags
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
				tags:         instance.Tags,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

func (f *LightsailInstance) FeatureFlags(ff config.FeatureFlags) {
	f.featureFlags = ff
}

func (f *LightsailInstance) Remove() error {

	_, err := f.svc.DeleteInstance(&lightsail.DeleteInstanceInput{
		InstanceName:      f.instanceName,
		ForceDeleteAddOns: aws.Bool(f.featureFlags.ForceDeleteLightsailAddOns),
	})

	return err
}

func (f *LightsailInstance) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.Set("Name", f.instanceName)
	return properties
}

func (f *LightsailInstance) String() string {
	return *f.instanceName
}
