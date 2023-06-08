package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gamelift"
)

type GameLiftQueue struct {
	svc  *gamelift.GameLift
	Name string
}

func init() {
	register("GameLiftQueue", ListGameLiftQueues)
}

func ListGameLiftQueues(sess *session.Session) ([]Resource, error) {
	svc := gamelift.New(sess)

	resp, err := svc.DescribeGameSessionQueues(&gamelift.DescribeGameSessionQueuesInput{})
	if err != nil {
		return nil, err
	}

	queues := make([]Resource, 0)
	for _, queue := range resp.GameSessionQueues {
		q := &GameLiftQueue{
			svc:  svc,
			Name: *queue.Name,
		}
		queues = append(queues, q)
	}

	return queues, nil
}

func (queue *GameLiftQueue) Remove() error {
	params := &gamelift.DeleteGameSessionQueueInput{
		Name: aws.String(queue.Name),
	}

	_, err := queue.svc.DeleteGameSessionQueue(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *GameLiftQueue) String() string {
	return i.Name
}
