package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMSKey struct {
	svc   *kms.KMS
	id    string
	state string
	alias string
}

func (n *KMSNuke) ListKeys() ([]Resource, error) {
	respAlias, err := n.Service.ListAliases(nil)
	if err != nil {
		return nil, err
	}

	aliasMap := map[string]string{}
	for _, alias := range respAlias.Aliases {
		if alias.TargetKeyId != nil {
			aliasMap[*alias.TargetKeyId] = *alias.AliasName
		}
	}

	resp, err := n.Service.ListKeys(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, key := range resp.Keys {
		resp, err := n.Service.DescribeKey(&kms.DescribeKeyInput{
			KeyId: key.KeyId,
		})
		if err != nil {
			return nil, err
		}

		resources = append(resources, &KMSKey{
			svc:   n.Service,
			id:    *resp.KeyMetadata.KeyId,
			state: *resp.KeyMetadata.KeyState,
			alias: aliasMap[*resp.KeyMetadata.KeyId],
		})
	}

	return resources, nil
}

func (e *KMSKey) Filter() error {
	if e.state == "PendingDeletion" {
		return fmt.Errorf("is already in PendingDeletion state")
	}

	if strings.HasPrefix(e.alias, "alias/aws/") {
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
