package resources

import "github.com/aws/aws-sdk-go/service/iam"

type IAMServerCertificate struct {
	svc  *iam.IAM
	name string
}

func (n *IAMNuke) ListServerCertificates() ([]Resource, error) {
	resp, err := n.Service.ListServerCertificates(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, meta := range resp.ServerCertificateMetadataList {
		resources = append(resources, &IAMServerCertificate{
			svc:  n.Service,
			name: *meta.ServerCertificateName,
		})
	}

	return resources, nil
}

func (e *IAMServerCertificate) Remove() error {
	_, err := e.svc.DeleteServerCertificate(&iam.DeleteServerCertificateInput{
		ServerCertificateName: &e.name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMServerCertificate) String() string {
	return e.name
}
