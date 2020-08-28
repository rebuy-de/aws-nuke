package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type KMSKey struct {
	svc     *kms.KMS
	id      string
	state   string
	manager *string
	tags    []*kms.Tag
}

func init() {
	register("KMSKey", ListKMSKeys)
}

func ListKMSKeys(sess *session.Session) ([]Resource, error) {
	svc := kms.New(sess)

	resp, err := svc.ListKeys(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, key := range resp.Keys {
		tags, err := svc.ListResourceTags(&kms.ListResourceTagsInput{
			KeyId: key.KeyId,
		})

		resp, err := svc.DescribeKey(&kms.DescribeKeyInput{
			KeyId: key.KeyId,
		})
		if err != nil {
			return nil, err
		}

		resources = append(resources, &KMSKey{
			svc:     svc,
			id:      *resp.KeyMetadata.KeyId,
			state:   *resp.KeyMetadata.KeyState,
			manager: resp.KeyMetadata.KeyManager,
			tags:    tags.Tags,
		})
	}

	return resources, nil
}

func (e *KMSKey) Filter() error {
	if e.state == "PendingDeletion" {
		return fmt.Errorf("is already in PendingDeletion state")
	}

	if e.manager != nil && *e.manager == kms.KeyManagerTypeAws {
		return fmt.Errorf("cannot delete AWS managed key")
	}

	return nil
}

func (e *KMSKey) Remove() error {
	_, err := e.svc.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
		KeyId:               &e.id,
		PendingWindowInDays: aws.Int64(7),
	})
	return err
}

func (e *KMSKey) String() string {
	return e.id
}

func (i *KMSKey) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("ID", i.id)

	for _, tag := range i.tags {
		properties.SetTag(tag.TagKey, tag.TagValue)
	}

	return properties
}
