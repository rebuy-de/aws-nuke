package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type LightsailDisk struct {
	svc      *lightsail.Lightsail
	diskName *string
}

func init() {
	register("LightsailDisk", ListLightsailDisks)
}

func ListLightsailDisks(sess *session.Session) ([]Resource, error) {
	svc := lightsail.New(sess)
	resources := []Resource{}

	params := &lightsail.GetDisksInput{}

	for {
		output, err := svc.GetDisks(params)
		if err != nil {
			return nil, err
		}

		for _, disk := range output.Disks {
			resources = append(resources, &LightsailDisk{
				svc:      svc,
				diskName: disk.Name,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

func (f *LightsailDisk) Remove() error {

	_, err := f.svc.DeleteDisk(&lightsail.DeleteDiskInput{
		DiskName: f.diskName,
	})

	return err
}

func (f *LightsailDisk) String() string {
	return *f.diskName
}
