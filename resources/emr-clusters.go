package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
)

type EMRCluster struct {
	svc   *emr.EMR
	ID    *string
	state *string
}

func init() {
	register("EMRCluster", ListEMRClusters)
}

func ListEMRClusters(sess *session.Session) ([]Resource, error) {
	svc := emr.New(sess)
	resources := []Resource{}

	params := &emr.ListClustersInput{}

	for {
		resp, err := svc.ListClusters(params)
		if err != nil {
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			resources = append(resources, &EMRCluster{
				svc:   svc,
				ID:    cluster.Id,
				state: cluster.Status.State,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (f *EMRCluster) Remove() error {

	//Call names are inconsistent in the SDK
	_, err := f.svc.TerminateJobFlows(&emr.TerminateJobFlowsInput{
		JobFlowIds: []*string{f.ID},
	})
	// Force nil return due to async callbacks blocking
	if err == nil {
		return nil
	}

	return err
}

func (f *EMRCluster) String() string {
	return *f.ID
}

func (f *EMRCluster) Filter() error {
	if strings.Contains(*f.state, "TERMINATED") {
		return fmt.Errorf("already terminated")
	}
	return nil
}
