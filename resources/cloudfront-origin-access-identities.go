package resources

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontOriginAccessIdentity struct {
    svc *cloudfront.CloudFront
    ID  *string
}

func init() {
    register("CloudFrontOriginAccessIdentity", ListCloudFrontOriginAccessIdentities)
}

func ListCloudFrontOriginAccessIdentities(sess *session.Session) ([]Resource, error) {
  svc := cloudfront.New(sess)
  resources := []Resource{}

  for {
    resp, err := svc.ListCloudFrontOriginAccessIdentities(nil)
    if err != nil {
        return nil, err
    }

    for _, item := range resp.CloudFrontOriginAccessIdentityList.Items {
      resources = append(resources,&CloudFrontOriginAccessIdentity{
        svc: svc,
        ID:  item.Id,
      })
    }
    return resources, nil
  }
}

func (f *CloudFrontOriginAccessIdentity) Remove() error {
  resp, err := f.svc.GetCloudFrontOriginAccessIdentity(&cloudfront.GetCloudFrontOriginAccessIdentityInput{
    Id: f.ID,
  })
  if err != nil {
    return err
  }

  _, err = f.svc.DeleteCloudFrontOriginAccessIdentity(&cloudfront.DeleteCloudFrontOriginAccessIdentityInput{
    Id: f.ID,
    IfMatch: resp.ETag,
  })

  return err
}

func (f *CloudFrontOriginAccessIdentity) Properties() types.Properties {
  properties := types.NewProperties()
  properties.Set("ID", f.ID)
  return properties
}
