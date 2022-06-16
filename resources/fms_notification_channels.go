package resources

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/fms"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
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
			if aerr.Code() != fms.ErrCodeResourceNotFoundException {
				return nil, err
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
