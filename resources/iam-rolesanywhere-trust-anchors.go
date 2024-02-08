package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rolesanywhere"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TrustAnchor struct {
	svc  *rolesanywhere.RolesAnywhere
	TrustAnchorId   string
}

func init() {
	register("IAMRolesAnywhereTrustAnchors", ListTrustAnchors)
}

func ListTrustAnchors(sess *session.Session) ([]Resource, error) {
	svc := rolesanywhere.New(sess)

	params := &rolesanywhere.ListTrustAnchorsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListTrustAnchors(params)
		if err != nil {
			return nil, err
		}
		for _, trustAnchor := range resp.TrustAnchors {
			resources = append(resources, &TrustAnchor{
				svc:      svc,
				TrustAnchorId: *trustAnchor.TrustAnchorId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (e *TrustAnchor) Remove() error {
	_, err := e.svc.DeleteTrustAnchor(&rolesanywhere.DeleteTrustAnchorInput{
		TrustAnchorId: &e.TrustAnchorId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *TrustAnchor) String() string {
	return e.TrustAnchorId
}

func (e *TrustAnchor) Properties() types.Properties {
	return types.NewProperties().
		Set("TrustAnchorId", e.TrustAnchorId)
}
