package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type NetworkManagerCoreNetwork struct {
	svc     *networkmanager.NetworkManager
	network *networkmanager.CoreNetworkSummary
}

func init() {
	register("NetworkManagerCoreNetwork", ListNetworkManagerCoreNetworks)
}

func ListNetworkManagerCoreNetworks(sess *session.Session) ([]Resource, error) {
	svc := networkmanager.New(sess)
	params := &networkmanager.ListCoreNetworksInput{}
	resources := make([]Resource, 0)

	resp, err := svc.ListCoreNetworks(params)
	if err != nil {
		return nil, err
	}

	for _, network := range resp.CoreNetworks {
		resources = append(resources, &NetworkManagerCoreNetwork{
			svc:     svc,
			network: network,
		})
	}

	return resources, nil
}

func (n *NetworkManagerCoreNetwork) Remove() error {
	params := &networkmanager.DeleteCoreNetworkInput{
		CoreNetworkId: n.network.CoreNetworkId,
	}

	_, err := n.svc.DeleteCoreNetwork(params)
	if err != nil {
		return err
	}

	return nil

}

func (n *NetworkManagerCoreNetwork) Filter() error {
	if strings.ToLower(*n.network.State) == "deleted" {
		return fmt.Errorf("already deleted")
	}

	return nil
}

func (n *NetworkManagerCoreNetwork) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range n.network.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.
		Set("ID", n.network.CoreNetworkId).
		Set("ARN", n.network.CoreNetworkArn)
	return properties
}

func (n *NetworkManagerCoreNetwork) String() string {
	return *n.network.CoreNetworkId
}
