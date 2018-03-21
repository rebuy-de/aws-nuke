package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediapackage"
)

type MediaPackageOriginEndpoint struct {
	svc *mediapackage.MediaPackage
	ID  *string
}

func init() {
	register("MediaPackageOriginEndpoint", ListMediaPackageOriginEndpoints)
}

func ListMediaPackageOriginEndpoints(sess *session.Session) ([]Resource, error) {
	svc := mediapackage.New(sess)
	resources := []Resource{}

	params := &mediapackage.ListOriginEndpointsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListOriginEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, originEndpoint := range output.OriginEndpoints {
			resources = append(resources, &MediaPackageOriginEndpoint{
				svc: svc,
				ID:  originEndpoint.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaPackageOriginEndpoint) Remove() error {

	_, err := f.svc.DeleteOriginEndpoint(&mediapackage.DeleteOriginEndpointInput{
		Id: f.ID,
	})

	return err
}

func (f *MediaPackageOriginEndpoint) String() string {
	return *f.ID
}
