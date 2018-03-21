package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTCACertificate struct {
	svc *iot.IoT
	ID  *string
}

func init() {
	register("IoTCACertificate", ListIoTCACertificates)
}

func ListIoTCACertificates(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListCACertificatesInput{}

	output, err := svc.ListCACertificates(params)
	if err != nil {
		return nil, err
	}

	for _, certificate := range output.Certificates {
		resources = append(resources, &IoTCACertificate{
			svc: svc,
			ID:  certificate.CertificateId,
		})
	}

	return resources, nil
}

func (f *IoTCACertificate) Remove() error {

	_, err := f.svc.UpdateCACertificate(&iot.UpdateCACertificateInput{
		CertificateId: f.ID,
		NewStatus:     aws.String("INACTIVE"),
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteCACertificate(&iot.DeleteCACertificateInput{
		CertificateId: f.ID,
	})

	return err
}

func (f *IoTCACertificate) String() string {
	return *f.ID
}
