package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type SageMakerDomain struct {
	svc      *sagemaker.SageMaker
	domainID *string
}

func init() {
	register("SageMakerDomain", ListSageMakerDomains)
}

func ListSageMakerDomains(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListDomainsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListDomains(params)
		if err != nil {
			return nil, err
		}

		for _, domain := range resp.Domains {
			resources = append(resources, &SageMakerDomain{
				svc:      svc,
				domainID: domain.DomainId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerDomain) Remove() error {
	_, err := f.svc.DeleteDomain(&sagemaker.DeleteDomainInput{
		DomainId:        f.domainID,
		RetentionPolicy: &sagemaker.RetentionPolicy{HomeEfsFileSystem: aws.String(sagemaker.RetentionTypeDelete)},
	})

	return err
}

func (f *SageMakerDomain) String() string {
	return *f.domainID
}

func (i *SageMakerDomain) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("DomainID", i.domainID)
	return properties
}
