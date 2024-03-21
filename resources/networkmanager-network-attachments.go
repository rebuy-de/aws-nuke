package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type NetworkManagerNetworkAttachment struct {
	svc        *networkmanager.NetworkManager
	attachment *networkmanager.Attachment
}

func init() {
	register("NetworkManagerNetworkAttachment", ListNetworkManagerNetworkAttachments)
}

func ListNetworkManagerNetworkAttachments(sess *session.Session) ([]Resource, error) {
	svc := networkmanager.New(sess)
	params := &networkmanager.ListAttachmentsInput{}
	resources := make([]Resource, 0)

	resp, err := svc.ListAttachments(params)
	if err != nil {
		return nil, err
	}

	for _, attachment := range resp.Attachments {
		resources = append(resources, &NetworkManagerNetworkAttachment{
			svc:        svc,
			attachment: attachment,
		})
	}

	return resources, nil
}

func (n *NetworkManagerNetworkAttachment) Remove() error {
	params := &networkmanager.DeleteAttachmentInput{
		AttachmentId: n.attachment.AttachmentId,
	}

	_, err := n.svc.DeleteAttachment(params)
	if err != nil {
		return err
	}

	return nil

}

func (n *NetworkManagerNetworkAttachment) Filter() error {
	if strings.ToLower(*n.attachment.State) == "deleted" {
		return fmt.Errorf("already deleted")
	}

	return nil
}

func (n *NetworkManagerNetworkAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range n.attachment.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.
		Set("ID", n.attachment.AttachmentId).
		Set("ARN", n.attachment.ResourceArn)
	return properties
}

func (n *NetworkManagerNetworkAttachment) String() string {
	return *n.attachment.AttachmentId
}
