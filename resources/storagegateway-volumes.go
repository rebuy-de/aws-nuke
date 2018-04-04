package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/storagegateway"
)

type StorageGatewayVolume struct {
	svc *storagegateway.StorageGateway
	ARN *string
}

func init() {
	register("StorageGatewayVolume", ListStorageGatewayVolumes)
}

func ListStorageGatewayVolumes(sess *session.Session) ([]Resource, error) {
	svc := storagegateway.New(sess)
	resources := []Resource{}

	params := &storagegateway.ListVolumesInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListVolumes(params)
		if err != nil {
			return nil, err
		}

		for _, volumeInfo := range output.VolumeInfos {
			resources = append(resources, &StorageGatewayVolume{
				svc: svc,
				ARN: volumeInfo.VolumeARN,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *StorageGatewayVolume) Remove() error {

	_, err := f.svc.DeleteVolume(&storagegateway.DeleteVolumeInput{
		VolumeARN: f.ARN,
	})

	return err
}

func (f *StorageGatewayVolume) String() string {
	return *f.ARN
}
