package resources

import "github.com/aws/aws-sdk-go/service/iam"

type IamServerCertificate struct {
	svc  *iam.IAM
	name string
}

func (n *IamNuke) ListServerCertificates() ([]Resource, error) {
	resp, err := n.Service.ListServerCertificates(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, meta := range resp.ServerCertificateMetadataList {
		resources = append(resources, &IamServerCertificate{
			svc:  n.Service,
			name: *meta.ServerCertificateName,
		})
	}

	return resources, nil
}

func (e *IamServerCertificate) Remove() error {
	_, err := e.svc.DeleteServerCertificate(&iam.DeleteServerCertificateInput{
		ServerCertificateName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamServerCertificate) String() string {
	return e.name
}
