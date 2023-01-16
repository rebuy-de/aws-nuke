package resources

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EMRCluster struct {
	svc     *emr.EMR
	cluster *emr.ClusterSummary
	state   *string
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
				svc:     svc,
				cluster: cluster,
				state:   cluster.Status.State,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (f *EMRCluster) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreatedTime", f.cluster.Status.Timeline.CreationDateTime.Format(time.RFC3339))
	return properties
}

func (f *EMRCluster) Remove() error {

	//Call names are inconsistent in the SDK
	_, err := f.svc.TerminateJobFlows(&emr.TerminateJobFlowsInput{
		JobFlowIds: []*string{f.cluster.Id},
	})
	// Force nil return due to async callbacks blocking
	if err == nil {
		return nil
	}

	return err
}

func (f *EMRCluster) String() string {
	return *f.cluster.Id
}

func (f *EMRCluster) Filter() error {
	if strings.Contains(*f.state, "TERMINATED") {
		return fmt.Errorf("already terminated")
	}
	return nil
}
