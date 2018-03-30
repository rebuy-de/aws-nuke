package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworkscm"
)

type OpsWorksCMServer struct {
	svc    *opsworkscm.OpsWorksCM
	name   *string
	status *string
}

func init() {
	register("OpsWorksCMServer", ListOpsWorksCMServers)
}

func ListOpsWorksCMServers(sess *session.Session) ([]Resource, error) {
	svc := opsworkscm.New(sess)
	resources := []Resource{}

	params := &opsworkscm.DescribeServersInput{}

	output, err := svc.DescribeServers(params)
	if err != nil {
		return nil, err
	}

	for _, server := range output.Servers {
		resources = append(resources, &OpsWorksCMServer{
			svc:  svc,
			name: server.ServerName,
		})
	}

	return resources, nil
}

func (f *OpsWorksCMServer) Remove() error {

	_, err := f.svc.DeleteServer(&opsworkscm.DeleteServerInput{
		ServerName: f.name,
	})

	return err
}

func (f *OpsWorksCMServer) String() string {
	return *f.name
}
