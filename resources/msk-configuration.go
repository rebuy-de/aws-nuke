package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MSKConfiguration struct {
	svc  *kafka.Kafka
	arn  string
	name string
}

func init() {
	register("MSKConfiguration", ListMSKConfigurations)
}

func ListMSKConfigurations(sess *session.Session) ([]Resource, error) {
	svc := kafka.New(sess)
	params := &kafka.ListConfigurationsInput{}
	resp, err := svc.ListConfigurations(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, configuration := range resp.Configurations {
		resources = append(resources, &MSKConfiguration{
			svc:  svc,
			arn:  *configuration.Arn,
			name: *configuration.Name,
		})
	}

	return resources, nil
}

func (m *MSKConfiguration) Remove() error {
	params := &kafka.DeleteConfigurationInput{
		Arn: &m.arn,
	}

	_, err := m.svc.DeleteConfiguration(params)
	if err != nil {
		return err
	}

	return nil
}

func (m *MSKConfiguration) String() string {
	return m.arn
}

func (m *MSKConfiguration) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ARN", m.arn)
	properties.Set("Name", m.name)

	return properties
}
