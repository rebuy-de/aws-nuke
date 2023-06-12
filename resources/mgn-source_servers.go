package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mgn"
)

type MgnSourceServer struct {
	svc *mgn.Mgn
	id  *string
}

func init() {
	register("MgnSourceServer", ListMgnSourceServers)
}

func ListMgnSourceServers(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &MgnSourceServer{
				svc: svc,
				id:  sourceServer.SourceServerID,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MgnSourceServer) Remove() error {

	_, err := f.svc.DeleteSourceServer(&mgn.DeleteSourceServerInput{
		SourceServerID: f.id,
	})

	return err
}

func (f *MgnSourceServer) String() string {
	return *f.id
}
