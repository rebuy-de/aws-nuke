package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transfer"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TransferServer struct {
	svc          *transfer.Transfer
	serverID     *string
	endpointType *string
	protocols    []string
	tags         []*transfer.Tag
}

func init() {
	register("TransferServer", ListTransferServers)
}

func ListTransferServers(sess *session.Session) ([]Resource, error) {
	svc := transfer.New(sess)
	resources := []Resource{}

	params := &transfer.ListServersInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListServers(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Servers {
			descOutput, err := svc.DescribeServer(&transfer.DescribeServerInput{
				ServerId: item.ServerId,
			})
			if err != nil {
				return nil, err
			}

			protocols := []string{}
			for _, protocol := range descOutput.Server.Protocols {
				protocols = append(protocols, *protocol)
			}

			resources = append(resources, &TransferServer{
				svc:          svc,
				serverID:     item.ServerId,
				endpointType: item.EndpointType,
				protocols:    protocols,
				tags:         descOutput.Server.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (ts *TransferServer) Remove() error {

	_, err := ts.svc.DeleteServer(&transfer.DeleteServerInput{
		ServerId: ts.serverID,
	})

	return err
}

func (ts *TransferServer) String() string {
	return *ts.serverID
}

func (ts *TransferServer) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range ts.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.
		Set("ServerID", ts.serverID).
		Set("EndpointType", ts.endpointType).
		Set("Protocols", strings.Join(ts.protocols, ", "))
	return properties
}
