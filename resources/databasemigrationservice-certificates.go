package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

type DatabaseMigrationServiceCertificate struct {
	svc *databasemigrationservice.DatabaseMigrationService
	ARN *string
}

func init() {
	register("DatabaseMigrationServiceCertificate", ListDatabaseMigrationServiceCertificates)
}

func ListDatabaseMigrationServiceCertificates(sess *session.Session) ([]Resource, error) {
	svc := databasemigrationservice.New(sess)
	resources := []Resource{}

	params := &databasemigrationservice.DescribeCertificatesInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeCertificates(params)
		if err != nil {
			return nil, err
		}

		for _, certificate := range output.Certificates {
			resources = append(resources, &DatabaseMigrationServiceCertificate{
				svc: svc,
				ARN: certificate.CertificateArn,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *DatabaseMigrationServiceCertificate) Remove() error {

	_, err := f.svc.DeleteEndpoint(&databasemigrationservice.DeleteEndpointInput{
		EndpointArn: f.ARN,
	})

	return err
}

func (f *DatabaseMigrationServiceCertificate) String() string {
	return *f.ARN
}
