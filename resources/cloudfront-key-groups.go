package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontKeyGroup struct {
	svc              *cloudfront.CloudFront
	ID               *string
	name             *string
	lastModifiedTime *time.Time
}

func init() {
	register("CloudFrontKeyGroup", ListCloudFrontKeyGroups)
}

func ListCloudFrontKeyGroups(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}
	params := &cloudfront.ListKeyGroupsInput{}

	for {
		resp, err := svc.ListKeyGroups(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.KeyGroupList.Items {
			resources = append(resources, &CloudFrontKeyGroup{
				svc:              svc,
				ID:               item.KeyGroup.Id,
				name:             item.KeyGroup.KeyGroupConfig.Name,
				lastModifiedTime: item.KeyGroup.LastModifiedTime,
			})
		}

		if resp.KeyGroupList.NextMarker == nil {
			break
		}

		params.Marker = resp.KeyGroupList.NextMarker
	}

	return resources, nil
}

func (f *CloudFrontKeyGroup) Remove() error {
	resp, err := f.svc.GetKeyGroup(&cloudfront.GetKeyGroupInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteKeyGroup(&cloudfront.DeleteKeyGroupInput{
		Id:      f.ID,
		IfMatch: resp.ETag,
	})

	return err
}

func (f *CloudFrontKeyGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	properties.Set("Name", f.name)
	properties.Set("LastModifiedTime", f.lastModifiedTime.Format(time.RFC3339))
	return properties
}
