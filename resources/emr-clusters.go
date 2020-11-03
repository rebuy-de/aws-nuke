package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/emr/emriface"
	"github.com/rebuy-de/aws-nuke/pkg/config"
)

type EMRCluster struct {
	svc   emriface.EMRAPI
	ID    *string
	state *string

	featureFlags config.FeatureFlags
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

func (c *EMRCluster) FeatureFlags(ff config.FeatureFlags) {
	c.featureFlags = ff
}

func (f *EMRCluster) Remove() error {
	params := &emr.TerminateJobFlowsInput{
		JobFlowIds: []*string{f.ID},
	}

	//Call names are inconsistent in the SDK
	_, err := f.svc.TerminateJobFlows(params)

	if err != nil {
		if f.featureFlags.DisableDeletionProtection.EMRCluster {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "ValidationException" &&
				awsErr.Message() == "Could not shut down one or more job flows since they are termination protected." {
				err = f.DisableProtection()
				if err != nil {
					return err
				}
				_, err := f.svc.TerminateJobFlows(params)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return err
	}

	return err
}

func (c *EMRCluster) DisableProtection() error {
	params := &emr.SetTerminationProtectionInput{
		JobFlowIds:           []*string{c.ID},
		TerminationProtected: aws.Bool(false),
	}
	_, err := c.svc.SetTerminationProtection(params)
	if err != nil {
		return err
	}
	return nil
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
