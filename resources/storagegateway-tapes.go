package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/storagegateway"
)

type StorageGatewayTape struct {
	svc        *storagegateway.StorageGateway
	tapeARN    *string
	gatewayARN *string
}

func init() {
	register("StorageGatewayTape", ListStorageGatewayTapes)
}

func ListStorageGatewayTapes(sess *session.Session) ([]Resource, error) {
	svc := storagegateway.New(sess)
	resources := []Resource{}

	params := &storagegateway.ListTapesInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListTapes(params)
		if err != nil {
			return nil, err
		}

		for _, tapeInfo := range output.TapeInfos {
			resources = append(resources, &StorageGatewayTape{
				svc:        svc,
				tapeARN:    tapeInfo.TapeARN,
				gatewayARN: tapeInfo.GatewayARN,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *StorageGatewayTape) Remove() error {

	_, err := f.svc.DeleteTape(&storagegateway.DeleteTapeInput{
		TapeARN:    f.tapeARN,
		GatewayARN: f.gatewayARN,
	})

	return err
}

func (f *StorageGatewayTape) String() string {
	return *f.tapeARN
}
