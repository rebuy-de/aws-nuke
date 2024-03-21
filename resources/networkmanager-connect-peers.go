package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type NetworkManagerConnectPeer struct {
	svc  *networkmanager.NetworkManager
	peer *networkmanager.ConnectPeerSummary
}

func init() {
	register("NetworkManagerConnectPeer", ListNetworkManagerConnectPeers)
}

func ListNetworkManagerConnectPeers(sess *session.Session) ([]Resource, error) {
	svc := networkmanager.New(sess)
	params := &networkmanager.ListConnectPeersInput{}
	resources := make([]Resource, 0)

	resp, err := svc.ListConnectPeers(params)
	if err != nil {
		return nil, err
	}

	for _, connectPeer := range resp.ConnectPeers {
		resources = append(resources, &NetworkManagerConnectPeer{
			svc:  svc,
			peer: connectPeer,
		})
	}

	return resources, nil
}

func (n *NetworkManagerConnectPeer) Remove() error {
	params := &networkmanager.DeleteConnectPeerInput{
		ConnectPeerId: n.peer.ConnectPeerId,
	}

	_, err := n.svc.DeleteConnectPeer(params)
	if err != nil {
		return err
	}

	return nil

}

func (n *NetworkManagerConnectPeer) Filter() error {
	if strings.ToLower(*n.peer.ConnectPeerState) == "deleted" {
		return fmt.Errorf("already deleted")
	}

	return nil
}

func (n *NetworkManagerConnectPeer) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range n.peer.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("ID", n.peer.ConnectPeerId)
	return properties
}

func (n *NetworkManagerConnectPeer) String() string {
	return *n.peer.ConnectPeerId
}
