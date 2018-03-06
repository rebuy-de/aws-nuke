package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
)

type ConfigServiceDeliveryChannel struct {
	svc                 *configservice.ConfigService
	deliveryChannelName *string
}

func init() {
	register("ConfigServiceDeliveryChannel", ListConfigServiceDeliveryChannels)
}

func ListConfigServiceDeliveryChannels(sess *session.Session) ([]Resource, error) {
	svc := configservice.New(sess)

	params := &configservice.DescribeDeliveryChannelsInput{}
	resp, err := svc.DescribeDeliveryChannels(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, deliveryChannel := range resp.DeliveryChannels {
		resources = append(resources, &ConfigServiceDeliveryChannel{
			svc:                 svc,
			deliveryChannelName: deliveryChannel.Name,
		})
	}

	return resources, nil
}

func (f *ConfigServiceDeliveryChannel) Remove() error {

	_, err := f.svc.DeleteDeliveryChannel(&configservice.DeleteDeliveryChannelInput{
		DeliveryChannelName: f.deliveryChannelName,
	})

	return err
}

func (f *ConfigServiceDeliveryChannel) String() string {
	return *f.deliveryChannelName
}
