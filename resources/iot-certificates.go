package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTCertificate struct {
	svc *iot.IoT
	ID  *string
}

func init() {
	register("IoTCertificate", ListIoTCertificates)
}

func ListIoTCertificates(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListCertificatesInput{}

	for {
		output, err := svc.ListCertificates(params)
		if err != nil {
			return nil, err
		}

		for _, certificate := range output.Certificates {
			resources = append(resources, &IoTCertificate{
				svc: svc,
				ID:  certificate.CertificateId,
			})
		}
		if output.NextMarker == nil {
			break
		}

		params.Marker = output.NextMarker
	}

	return resources, nil
}

func (f *IoTCertificate) Remove() error {

	_, err := f.svc.UpdateCertificate(&iot.UpdateCertificateInput{
		CertificateId: f.ID,
		NewStatus:     aws.String("INACTIVE"),
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteCertificate(&iot.DeleteCertificateInput{
		CertificateId: f.ID,
	})

	return err
}

func (f *IoTCertificate) String() string {
	return *f.ID
}
