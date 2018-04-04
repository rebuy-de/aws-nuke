package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/storagegateway"
)

type StorageGatewayFileShare struct {
	svc *storagegateway.StorageGateway
	ARN *string
}

func init() {
	register("StorageGatewayFileShare", ListStorageGatewayFileShares)
}

func ListStorageGatewayFileShares(sess *session.Session) ([]Resource, error) {
	svc := storagegateway.New(sess)
	resources := []Resource{}

	params := &storagegateway.ListFileSharesInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListFileShares(params)
		if err != nil {
			return nil, err
		}

		for _, fileShareInfo := range output.FileShareInfoList {
			resources = append(resources, &StorageGatewayFileShare{
				svc: svc,
				ARN: fileShareInfo.FileShareARN,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *StorageGatewayFileShare) Remove() error {

	_, err := f.svc.DeleteFileShare(&storagegateway.DeleteFileShareInput{
		FileShareARN: f.ARN,
		ForceDelete:  aws.Bool(true),
	})

	return err
}

func (f *StorageGatewayFileShare) String() string {
	return *f.ARN
}
