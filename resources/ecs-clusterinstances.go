package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ECSClusterInstance struct {
	svc         *ecs.ECS
	instanceARN *string
	clusterARN  *string
}

func init() {
	register("ECSClusterInstance", ListECSClusterInstances)
}

func ListECSClusterInstances(sess *session.Session) ([]Resource, error) {
	svc := ecs.New(sess)
	resources := []Resource{}
	clusters := []*string{}

	clusterParams := &ecs.ListClustersInput{
		MaxResults: aws.Int64(100),
	}

	// Iterate over clusters to ensure we dont presume its always default associations
	for {
		output, err := svc.ListClusters(clusterParams)
		if err != nil {
			return nil, err
		}

		for _, clusterArn := range output.ClusterArns {
			clusters = append(clusters, clusterArn)
		}

		if output.NextToken == nil {
			break
		}

		clusterParams.NextToken = output.NextToken
	}

	// Iterate over known clusters and discover their instances
	// to prevent assuming default is always used
	for _, clusterArn := range clusters {
		instanceParams := &ecs.ListContainerInstancesInput{
			Cluster:    clusterArn,
			MaxResults: aws.Int64(100),
		}
		output, err := svc.ListContainerInstances(instanceParams)
		if err != nil {
			return nil, err
		}

		for _, instanceArn := range output.ContainerInstanceArns {
			resources = append(resources, &ECSClusterInstance{
				svc:         svc,
				instanceARN: instanceArn,
				clusterARN:  clusterArn,
			})
		}

		if output.NextToken == nil {
			break
		}

		instanceParams.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ECSClusterInstance) Remove() error {

	_, err := f.svc.DeregisterContainerInstance(&ecs.DeregisterContainerInstanceInput{
		Cluster:           f.clusterARN,
		ContainerInstance: f.instanceARN,
		Force:             aws.Bool(true),
	})

	return err
}

func (f *ECSClusterInstance) String() string {
	return fmt.Sprintf("%s -> %s", *f.instanceARN, *f.clusterARN)
}
