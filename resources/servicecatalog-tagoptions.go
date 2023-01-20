package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

type ServiceCatalogTagOption struct {
	svc   *servicecatalog.ServiceCatalog
	ID    *string
	key   *string
	value *string
}

func init() {
	register("ServiceCatalogTagOption", ListServiceCatalogTagOptions)
}

func ListServiceCatalogTagOptions(sess *session.Session) ([]Resource, error) {
	svc := servicecatalog.New(sess)
	resources := []Resource{}

	params := &servicecatalog.ListTagOptionsInput{
		PageSize: aws.Int64(20),
	}

	for {
		resp, err := svc.ListTagOptions(params)
		if err != nil {
			if IsAWSError(err, servicecatalog.ErrCodeTagOptionNotMigratedException) {
				log.Info(err)
				break
			}
			return nil, err
		}

		for _, tagOptionDetail := range resp.TagOptionDetails {
			resources = append(resources, &ServiceCatalogTagOption{
				svc:   svc,
				ID:    tagOptionDetail.Id,
				key:   tagOptionDetail.Key,
				value: tagOptionDetail.Value,
			})
		}

		if resp.PageToken == nil {
			break
		}

		params.PageToken = resp.PageToken
	}

	return resources, nil
}

func (f *ServiceCatalogTagOption) Remove() error {

	_, err := f.svc.DeleteTagOption(&servicecatalog.DeleteTagOptionInput{
		Id: f.ID,
	})

	return err
}

func (f *ServiceCatalogTagOption) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	properties.Set("Key", f.key)
	properties.Set("Value", f.value)
	return properties
}

func (f *ServiceCatalogTagOption) String() string {
	return *f.ID
}
