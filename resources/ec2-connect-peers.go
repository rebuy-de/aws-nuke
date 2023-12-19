package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2TGWConnectPeer struct {
	svc  *ec2.EC2
	peer *ec2.TransitGatewayConnectPeer
}

func init() {
	register("EC2TGWConnectPeer", ListEC2TGWConnectPeer)
}

func ListEC2TGWConnectPeer(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)

	// filter should be set as deleted vpc connections are returned
	params := &ec2.DescribeTransitGatewayConnectPeersInput{}

	resp, err := svc.DescribeTransitGatewayConnectPeers(params)
	if err != nil {
		return nil, err
	}

	for _, connectPeer := range resp.TransitGatewayConnectPeers {
		resources = append(resources, &EC2TGWConnectPeer{
			svc:  svc,
			peer: connectPeer,
		})
	}

	return resources, nil
}

func (p *EC2TGWConnectPeer) Filter() error {
	if *p.peer.State == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (p *EC2TGWConnectPeer) Remove() error {
	params := &ec2.DeleteTransitGatewayConnectPeerInput{
		TransitGatewayConnectPeerId: p.peer.TransitGatewayConnectPeerId,
	}

	_, err := p.svc.DeleteTransitGatewayConnectPeer(params)
	if err != nil {
		return err
	}
	return nil
}

func (p *EC2TGWConnectPeer) String() string {
	return *p.peer.TransitGatewayConnectPeerId
}
