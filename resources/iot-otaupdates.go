package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTOTAUpdate struct {
	svc *iot.IoT
	ID  *string
}

func init() {
	register("IoTOTAUpdate", ListIoTOTAUpdates)
}

func ListIoTOTAUpdates(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListOTAUpdatesInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListOTAUpdates(params)
		if err != nil {
			return nil, err
		}

		for _, otaUpdate := range output.OtaUpdates {
			resources = append(resources, &IoTOTAUpdate{
				svc: svc,
				ID:  otaUpdate.OtaUpdateId,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTOTAUpdate) Remove() error {

	_, err := f.svc.DeleteOTAUpdate(&iot.DeleteOTAUpdateInput{
		OtaUpdateId: f.ID,
	})

	return err
}

func (f *IoTOTAUpdate) String() string {
	return *f.ID
}
