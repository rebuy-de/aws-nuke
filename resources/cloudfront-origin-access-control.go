package resources

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudfront"
  	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontOriginAccessControl struct {
    svc *cloudfront.CloudFront
    ID  *string
}

func init() {
    register("CloudFrontOriginAccessControl", ListCloudFrontOriginAccessControls)
}

func ListCloudFrontOriginAccessControls(sess *session.Session) ([]Resource, error) {
  svc := cloudfront.New(sess)
  resources := []Resource{}

  for {
    resp, err := svc.ListOriginAccessControls(nil)
    if err != nil {
        return nil, err
    }

    for _, item := range resp.OriginAccessControlList.Items {
      resources = append(resources,&CloudFrontOriginAccessControl{
        svc: svc,
        ID:  item.Id,
      })
    }
    return resources, nil
  }
}

func (f *CloudFrontOriginAccessControl) Remove() error {
  resp, err := f.svc.GetOriginAccessControl(&cloudfront.GetOriginAccessControlInput{
    Id: f.ID,
  })
  if err != nil {
    return err
  }

  _, err = f.svc.DeleteOriginAccessControl(&cloudfront.DeleteOriginAccessControlInput{
    Id: f.ID,
    IfMatch: resp.ETag,
  })

  return err
}

func (f *CloudFrontOriginAccessControl) Properties() types.Properties {
  properties := types.NewProperties()
  properties.Set("ID", f.ID)
  return properties
}
