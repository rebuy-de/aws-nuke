package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ECSTask struct {
	svc        *ecs.ECS
	taskARN    *string
	clusterARN *string
}

func init() {
	register("ECSTask", ListECSTasks)
}

func ListECSTasks(sess *session.Session) ([]Resource, error) {
	svc := ecs.New(sess)
	resources := []Resource{}
	clusters := []*string{}

	clusterParams := &ecs.ListClustersInput{
		MaxResults: aws.Int64(100),
	}

	// Discover all clusters
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

	// Discover all running tasks from all clusters
	for _, clusterArn := range clusters {
		taskParams := &ecs.ListTasksInput{
			Cluster:       clusterArn,
			MaxResults:    aws.Int64(10),
			DesiredStatus: aws.String("RUNNING"),
		}
		output, err := svc.ListTasks(taskParams)
		if err != nil {
			return nil, err
		}

		for _, taskArn := range output.TaskArns {
			resources = append(resources, &ECSTask{
				svc:        svc,
				taskARN:    taskArn,
				clusterARN: clusterArn,
			})
		}

		if output.NextToken == nil {
			continue
		}

		taskParams.NextToken = output.NextToken
	}

	return resources, nil
}

func (t *ECSTask) Filter() error {
	return nil
}

func (t *ECSTask) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("TaskARN", t.taskARN)
	properties.Set("ClusterARN", t.clusterARN)

	return properties
}

func (t *ECSTask) Remove() error {
	// When StopTask is called on a task, the equivalent of docker stop is issued to the
	// containers running in the task. This results in a SIGTERM value and a default
	// 30-second timeout, after which the SIGKILL value is sent and the containers are
	// forcibly stopped. If the container handles the SIGTERM value gracefully and exits
	// within 30 seconds from receiving it, no SIGKILL value is sent.

	_, err := t.svc.StopTask(&ecs.StopTaskInput{
		Cluster: t.clusterARN,
		Task:    t.taskARN,
		Reason:  aws.String("Task stopped via AWS Nuke"),
	})

	return err
}

func (t *ECSTask) String() string {
	return fmt.Sprintf("%s -> %s", *t.taskARN, *t.clusterARN)
}
