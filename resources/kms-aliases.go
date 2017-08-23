package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/kms"
)

type KMSAlias struct {
	svc  *kms.KMS
	name string
}

func (n *KMSNuke) ListAliases() ([]Resource, error) {
	resp, err := n.Service.ListAliases(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, alias := range resp.Aliases {
		resources = append(resources, &KMSAlias{
			svc:  n.Service,
			name: *alias.AliasName,
		})
	}

	return resources, nil
}

func (e *KMSAlias) Filter() error {
	if strings.HasPrefix(e.name, "alias/aws/") {
		return fmt.Errorf("cannot delete AWS alias")
	}
	return nil
}

func (e *KMSAlias) Remove() error {
	_, err := e.svc.DeleteAlias(&kms.DeleteAliasInput{
		AliasName: &e.name,
	})
	return err
}

func (e *KMSAlias) String() string {
	return e.name
}
