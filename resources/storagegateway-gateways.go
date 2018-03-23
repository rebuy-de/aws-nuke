package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/storagegateway"
)

type StorageGatewayGateway struct {
	svc *storagegateway.StorageGateway
	ARN *string
}

func init() {
	register("StorageGatewayGateway", ListStorageGatewayGateways)
}

func ListStorageGatewayGateways(sess *session.Session) ([]Resource, error) {
	svc := storagegateway.New(sess)
	resources := []Resource{}

	params := &storagegateway.ListGatewaysInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListGateways(params)
		if err != nil {
			return nil, err
		}

		for _, gateway := range output.Gateways {
			resources = append(resources, &StorageGatewayGateway{
				svc: svc,
				ARN: gateway.GatewayARN,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *StorageGatewayGateway) Remove() error {

	_, err := f.svc.DeleteGateway(&storagegateway.DeleteGatewayInput{
		GatewayARN: f.ARN,
	})

	return err
}

func (f *StorageGatewayGateway) String() string {
	return *f.ARN
}
