package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codestarconnections"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeStarConnection struct {
	svc              *codestarconnections.CodeStarConnections
	connectionARN    *string
	connectionName   *string
	providerType     *string
}

func init() {
	register("CodeStarConnection", ListCodeStarConnections)
}

func ListCodeStarConnections(sess *session.Session) ([]Resource, error) {
	svc := codestarconnections.New(sess)
	resources := []Resource{}

	params := &codestarconnections.ListConnectionsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListConnections(params)
		if err != nil {
			return nil, err
		}

		for _, connection := range output.Connections {
			resources = append(resources, &CodeStarConnection{
				svc:              svc,
				connectionARN:    connection.ConnectionArn,
				connectionName:   connection.ConnectionName,
				providerType:     connection.ProviderType,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CodeStarConnection) Remove() error {

	_, err := f.svc.DeleteConnection(&codestarconnections.DeleteConnectionInput{
		ConnectionArn: f.connectionARN,
	})

	return err
}

func (f *CodeStarConnection) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Name", f.connectionName).
		Set("ProviderType", f.providerType)
	return properties
}


func (f *CodeStarConnection) String() string {
	return *f.connectionName
}
