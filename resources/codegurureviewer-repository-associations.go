package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codegurureviewer"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeGuruReviewerRepositoryAssociation struct {
	svc            *codegurureviewer.CodeGuruReviewer
	AssociationArn *string
	AssociationId  *string
	Name           *string
	Owner          *string
	ProviderType   *string
}

func init() {
	register("CodeGuruReviewerRepositoryAssociation", ListCodeGuruReviewerRepositoryAssociations,
		mapCloudControl("AWS::CodeGuruReviewer::RepositoryAssociation"))
}

func ListCodeGuruReviewerRepositoryAssociations(sess *session.Session) ([]Resource, error) {
	svc := codegurureviewer.New(sess)
	resources := []Resource{}

	params := &codegurureviewer.ListRepositoryAssociationsInput{}

	for {
		resp, err := svc.ListRepositoryAssociations(params)
		if err != nil {
			return nil, err
		}

		for _, association := range resp.RepositoryAssociationSummaries {
			resources = append(resources, &CodeGuruReviewerRepositoryAssociation{
				svc:            svc,
				AssociationArn: association.AssociationArn,
				AssociationId:  association.AssociationId,
				Name:           association.Name,
				Owner:          association.Owner,
				ProviderType:   association.ProviderType,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeGuruReviewerRepositoryAssociation) Remove() error {
	_, err := f.svc.DisassociateRepository(&codegurureviewer.DisassociateRepositoryInput{
		AssociationArn: f.AssociationArn,
	})
	return err
}

func (f *CodeGuruReviewerRepositoryAssociation) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("AssociationArn", f.AssociationArn)
	properties.Set("AssociationId", f.AssociationId)
	properties.Set("Name", f.Name)
	properties.Set("Owner", f.Owner)
	properties.Set("ProviderType", f.ProviderType)
	return properties
}
