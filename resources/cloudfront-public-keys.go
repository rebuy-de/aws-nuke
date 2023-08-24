package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontPublicKey struct {
	svc         *cloudfront.CloudFront
	ID          *string
	name        *string
	createdTime *time.Time
}

func init() {
	register("CloudFrontPublicKey", ListCloudFrontPublicKeys)
}

func ListCloudFrontPublicKeys(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}
	params := &cloudfront.ListPublicKeysInput{}

	for {
		resp, err := svc.ListPublicKeys(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.PublicKeyList.Items {
			resources = append(resources, &CloudFrontPublicKey{
				svc:         svc,
				ID:          item.Id,
				name:        item.Name,
				createdTime: item.CreatedTime,
			})
		}

		if resp.PublicKeyList.NextMarker == nil {
			break
		}

		params.Marker = resp.PublicKeyList.NextMarker
	}

	return resources, nil
}

func (f *CloudFrontPublicKey) Remove() error {
	resp, err := f.svc.GetPublicKey(&cloudfront.GetPublicKeyInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeletePublicKey(&cloudfront.DeletePublicKeyInput{
		Id:      f.ID,
		IfMatch: resp.ETag,
	})

	return err
}

func (f *CloudFrontPublicKey) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	properties.Set("Name", f.name)
	properties.Set("CreatedTime", f.createdTime.Format(time.RFC3339))
	return properties
}
