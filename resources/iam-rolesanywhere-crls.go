package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rolesanywhere"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type Crl struct {
	svc  *rolesanywhere.RolesAnywhere
	CrlId   string
}

func init() {
	register("IAMRolesAnywhereCrls", ListCRLs)
}

func ListCRLs(sess *session.Session) ([]Resource, error) {
	svc := rolesanywhere.New(sess)

	params := &rolesanywhere.ListCrlsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListCrls(params)
		if err != nil {
			return nil, err
		}
		for _, crl := range resp.Crls {
			resources = append(resources, &Crl{
				svc:      svc,
				CrlId: *crl.CrlId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (e *Crl) Remove() error {
	_, err := e.svc.DeleteCrl(&rolesanywhere.DeleteCrlInput{
		CrlId: &e.CrlId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Crl) String() string {
	return e.CrlId
}

func (e *Crl) Properties() types.Properties {
	return types.NewProperties().
		Set("CrlId", e.CrlId)
}
