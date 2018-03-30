package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworkscm"
)

type OpsWorksCMServerState struct {
	svc    *opsworkscm.OpsWorksCM
	name   *string
	status *string
}

func init() {
	register("OpsWorksCMServerState", ListOpsWorksCMServerStates)
}

func ListOpsWorksCMServerStates(sess *session.Session) ([]Resource, error) {
	svc := opsworkscm.New(sess)
	resources := []Resource{}

	params := &opsworkscm.DescribeServersInput{}

	output, err := svc.DescribeServers(params)
	if err != nil {
		return nil, err
	}

	for _, server := range output.Servers {
		resources = append(resources, &OpsWorksCMServerState{
			svc:    svc,
			name:   server.ServerName,
			status: server.Status,
		})
	}

	return resources, nil
}

func (f *OpsWorksCMServerState) Remove() error {
	return nil
}

func (f *OpsWorksCMServerState) String() string {
	return *f.name
}

func (f *OpsWorksCMServerState) Filter() error {
	if *f.status == "CREATING" {
		return nil
	} else {
		return fmt.Errorf("available for transition")
	}
}
