package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/fms"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type FMSNotificationChannel struct {
	svc *fms.FMS
}

func init() {
	register("FMSNotificationChannel", ListFMSNotificationChannel)
}

func ListFMSNotificationChannel(sess *session.Session) ([]Resource, error) {
	svc := fms.New(sess)
	resources := []Resource{}

	if _, err := svc.GetNotificationChannel(&fms.GetNotificationChannelInput{}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if strings.Contains(aerr.Message(), "No default admin could be found") {
				logrus.Infof("FMSNotificationChannel: %s. Ignore if you haven't set it up.", aerr.Message())
				return nil, nil
			}
		} else {
			return nil, err
		}
	} else {
		resources = append(resources, &FMSNotificationChannel{
			svc: svc,
		})
	}

	return resources, nil
}

func (f *FMSNotificationChannel) Remove() error {

	_, err := f.svc.DeleteNotificationChannel(&fms.DeleteNotificationChannelInput{})

	return err
}

func (f *FMSNotificationChannel) String() string {
	return "fms-notification-channel"
}

func (f *FMSNotificationChannel) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("NotificationChannelEnabled", "true")
	return properties
}
