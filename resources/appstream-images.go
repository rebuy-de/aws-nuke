package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appstream"
)

type AppStreamImage struct {
	svc        *appstream.AppStream
	name       *string
	visibility *string
}

func init() {
	register("AppStreamImage", ListAppStreamImages)
}

func ListAppStreamImages(sess *session.Session) ([]Resource, error) {
	svc := appstream.New(sess)
	resources := []Resource{}

	params := &appstream.DescribeImagesInput{}

	output, err := svc.DescribeImages(params)
	if err != nil {
		return nil, err
	}

	for _, image := range output.Images {
		resources = append(resources, &AppStreamImage{
			svc:        svc,
			name:       image.Name,
			visibility: image.Visibility,
		})
	}

	return resources, nil
}

func (f *AppStreamImage) Remove() error {

	_, err := f.svc.DeleteImage(&appstream.DeleteImageInput{
		Name: f.name,
	})

	return err
}

func (f *AppStreamImage) String() string {
	return *f.name
}

func (f *AppStreamImage) Filter() error {
	if strings.ToUpper(*f.visibility) == "PUBLIC" {
		return fmt.Errorf("cannot delete public AWS images")
	}
	return nil
}
