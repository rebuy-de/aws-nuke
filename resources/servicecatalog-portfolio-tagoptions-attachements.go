package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

type ServiceCatalogTagOptionPortfolioAttachment struct {
	svc            *servicecatalog.ServiceCatalog
	tagOptionID    *string
	resourceID     *string
	tagOptionKey   *string
	tagOptionValue *string
	resourceName   *string
}

func init() {
	register("ServiceCatalogTagOptionPortfolioAttachment", ListServiceCatalogTagOptionPortfolioAttachments)
}

func ListServiceCatalogTagOptionPortfolioAttachments(sess *session.Session) ([]Resource, error) {
	svc := servicecatalog.New(sess)
	resources := []Resource{}
	tagOptions := []*servicecatalog.TagOptionDetail{}

	params := &servicecatalog.ListTagOptionsInput{
		PageSize: aws.Int64(20),
	}

	//List all Tag Options
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
			tagOptions = append(tagOptions, tagOptionDetail)
		}

		if resp.PageToken == nil {
			break
		}

		params.PageToken = resp.PageToken
	}

	resourceParams := &servicecatalog.ListResourcesForTagOptionInput{
		PageSize: aws.Int64(20),
	}

	for _, tagOption := range tagOptions {

		resourceParams.TagOptionId = tagOption.Id
		resp, err := svc.ListResourcesForTagOption(resourceParams)
		if err != nil {
			return nil, err
		}

		for _, resourceDetail := range resp.ResourceDetails {
			resources = append(resources, &ServiceCatalogTagOptionPortfolioAttachment{
				svc:            svc,
				tagOptionID:    tagOption.Id,
				resourceID:     resourceDetail.Id,
				resourceName:   resourceDetail.Name,
				tagOptionKey:   tagOption.Key,
				tagOptionValue: tagOption.Value,
			})
		}

		if resp.PageToken == nil {
			break
		}

		resourceParams.PageToken = resp.PageToken
	}

	return resources, nil
}

func (f *ServiceCatalogTagOptionPortfolioAttachment) Remove() error {

	_, err := f.svc.DisassociateTagOptionFromResource(&servicecatalog.DisassociateTagOptionFromResourceInput{
		TagOptionId: f.tagOptionID,
		ResourceId:  f.resourceID,
	})

	return err
}

func (f *ServiceCatalogTagOptionPortfolioAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("TagOptionID", f.tagOptionID)
	properties.Set("TagOptionKey", f.tagOptionKey)
	properties.Set("TagOptionValue", f.tagOptionValue)
	properties.Set("ResourceID", f.resourceID)
	properties.Set("ResourceName", f.resourceName)
	return properties
}

func (f *ServiceCatalogTagOptionPortfolioAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.tagOptionID, *f.resourceID)
}
