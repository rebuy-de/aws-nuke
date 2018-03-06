package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
)

type ConfigServiceConfigurationRecorder struct {
	svc                       *configservice.ConfigService
	configurationRecorderName *string
}

func init() {
	register("ConfigServiceConfigurationRecorder", ListConfigServiceConfigurationRecorders)
}

func ListConfigServiceConfigurationRecorders(sess *session.Session) ([]Resource, error) {
	svc := configservice.New(sess)

	params := &configservice.DescribeConfigurationRecordersInput{}
	resp, err := svc.DescribeConfigurationRecorders(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, configurationRecorder := range resp.ConfigurationRecorders {
		resources = append(resources, &ConfigServiceConfigurationRecorder{
			svc: svc,
			configurationRecorderName: configurationRecorder.Name,
		})
	}

	return resources, nil
}

func (f *ConfigServiceConfigurationRecorder) Remove() error {

	_, err := f.svc.DeleteConfigurationRecorder(&configservice.DeleteConfigurationRecorderInput{
		ConfigurationRecorderName: f.configurationRecorderName,
	})

	return err
}

func (f *ConfigServiceConfigurationRecorder) String() string {
	return *f.configurationRecorderName
}
