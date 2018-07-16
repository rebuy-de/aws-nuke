package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ECSService struct {
	svc        *ecs.ECS
	serviceARN *string
	clusterARN *string
}

func init() {
	register("ECSService", ListECSServices)
}

func ListECSServices(sess *session.Session) ([]Resource, error) {
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
		serviceParams := &ecs.ListServicesInput{
			Cluster:    clusterArn,
			MaxResults: aws.Int64(10),
		}
		output, err := svc.ListServices(serviceParams)
		if err != nil {
			return nil, err
		}

		for _, serviceArn := range output.ServiceArns {
			resources = append(resources, &ECSService{
				svc:        svc,
				serviceARN: serviceArn,
				clusterARN: clusterArn,
			})
		}

		if output.NextToken == nil {
			continue
		}

		serviceParams.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ECSService) Remove() error {

	_, err := f.svc.DeleteService(&ecs.DeleteServiceInput{
		Cluster: f.clusterARN,
		Service: f.serviceARN,
		Force:   aws.Bool(true),
	})

	return err
}

func (f *ECSService) String() string {
	return fmt.Sprintf("%s -> %s", *f.serviceARN, *f.clusterARN)
}
