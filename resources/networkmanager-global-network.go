package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type NetworkManagerGlobalNetwork struct {
	svc     *networkmanager.NetworkManager
	network *networkmanager.GlobalNetwork
}

func init() {
	register("NetworkManagerGlobalNetwork", ListNetworkManagerGlobalNetworks)
}

func ListNetworkManagerGlobalNetworks(sess *session.Session) ([]Resource, error) {
	svc := networkmanager.New(sess)
	resources := []Resource{}

	params := &networkmanager.DescribeGlobalNetworksInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.DescribeGlobalNetworks(params)
		if err != nil {
			return nil, err
		}

		for _, network := range resp.GlobalNetworks {
			resources = append(resources, &NetworkManagerGlobalNetwork{
				svc:     svc,
				network: network,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (n *NetworkManagerGlobalNetwork) Remove() error {
	params := &networkmanager.DeleteGlobalNetworkInput{
		GlobalNetworkId: n.network.GlobalNetworkId,
	}

	_, err := n.svc.DeleteGlobalNetwork(params)
	if err != nil {
		return err
	}

	return nil

}

func (n *NetworkManagerGlobalNetwork) Filter() error {
	if strings.ToLower(*n.network.State) == "deleted" {
		return fmt.Errorf("already deleted")
	}

	return nil
}

func (n *NetworkManagerGlobalNetwork) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range n.network.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.
		Set("ID", n.network.GlobalNetworkId).
		Set("ARN", n.network.GlobalNetworkArn)
	return properties
}

func (n *NetworkManagerGlobalNetwork) String() string {
	return *n.network.GlobalNetworkId
}
