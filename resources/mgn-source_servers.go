package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mgn"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MGNSourceServer struct {
	svc            *mgn.Mgn
	sourceServerID *string
	arn            *string
	tags           map[string]*string
}

func init() {
	register("MGNSourceServer", ListMGNSourceServers)
}

func ListMGNSourceServers(sess *session.Session) ([]Resource, error) {
	svc := mgn.New(sess)
	resources := []Resource{}

	params := &mgn.DescribeSourceServersInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeSourceServers(params)
		if err != nil {
			return nil, err
		}

		for _, sourceServer := range output.Items {
			resources = append(resources, &MGNSourceServer{
				svc:            svc,
				sourceServerID: sourceServer.SourceServerID,
				arn:            sourceServer.Arn,
				tags:           sourceServer.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MGNSourceServer) Remove() error {

	_, err := f.svc.DeleteSourceServer(&mgn.DeleteSourceServerInput{
		SourceServerID: f.sourceServerID,
	})

	return err
}

func (f *MGNSourceServer) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("SourceServerID", f.sourceServerID)
	properties.Set("Arn", f.arn)

	for key, val := range f.tags {
		properties.SetTag(&key, val)
	}
	return properties
}

func (f *MGNSourceServer) String() string {
	return *f.sourceServerID
}
