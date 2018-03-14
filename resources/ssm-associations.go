package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMAssociation struct {
	svc           *ssm.SSM
	associationID *string
	instanceID    *string
}

func init() {
	register("SSMAssociation", ListSSMAssociations)
}

func ListSSMAssociations(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	params := &ssm.ListAssociationsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListAssociations(params)
		if err != nil {
			return nil, err
		}

		for _, association := range output.Associations {
			resources = append(resources, &SSMAssociation{
				svc:           svc,
				associationID: association.AssociationId,
				instanceID:    association.InstanceId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMAssociation) Remove() error {

	_, err := f.svc.DeleteAssociation(&ssm.DeleteAssociationInput{
		AssociationId: f.associationID,
	})

	return err
}

func (f *SSMAssociation) String() string {
	return *f.associationID
}
