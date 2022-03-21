package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
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
	resources := make([]Resource, 0)

	var innerErr error
	err := svc.ListKeysPages(nil, func(resp *kms.ListKeysOutput, lastPage bool) bool {
		for _, key := range resp.Keys {
			resp, err := svc.DescribeKey(&kms.DescribeKeyInput{
				KeyId: key.KeyId,
			})
			if err != nil {
				innerErr = err
				return false
			}

			if *resp.KeyMetadata.KeyManager == kms.KeyManagerTypeAws {
				continue
			}

			if *resp.KeyMetadata.KeyState == kms.KeyStatePendingDeletion {
				continue
			}

			kmsKey := &KMSKey{
				svc:     svc,
				id:      *resp.KeyMetadata.KeyId,
				state:   *resp.KeyMetadata.KeyState,
				manager: resp.KeyMetadata.KeyManager,
			}

			tags, err := svc.ListResourceTags(&kms.ListResourceTagsInput{
				KeyId: key.KeyId,
			})
			if err != nil {
				innerErr = err
				return false
			}

			kmsKey.tags = tags.Tags
			resources = append(resources, kmsKey)
		}

		if lastPage {
			return false
		}

		return true
	})

	if err != nil {
		return nil, err
	}

	if innerErr != nil {
		return nil, err
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
