package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

func init() {
	register("SecurityHub", ListHubs)
}

func ListHubs(sess *session.Session) ([]Resource, error) {
	svc := securityhub.New(sess)

	resources := make([]Resource, 0)

	resp, err := svc.DescribeHub(nil)

	if err != nil {
		if IsAWSError(err, securityhub.ErrCodeInvalidAccessException) {
			// Security Hub is not enabled for this region
			return resources, nil
		}
		return nil, err
	}

	resources = append(resources, &Hub{
		svc: svc,
		id:  resp.HubArn,
	})
	return resources, nil
}

type Hub struct {
	svc *securityhub.SecurityHub
	id  *string
}

func (hub *Hub) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Arn", hub.id)
	return properties
}

func (hub *Hub) Remove() error {
	_, err := hub.svc.DisableSecurityHub(&securityhub.DisableSecurityHubInput{})
	return err
}
