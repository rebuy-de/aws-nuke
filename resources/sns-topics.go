package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sns"
)

func (n *SNSNuke) ListTopics() ([]Resource, error) {
	resp, err := n.Service.ListTopics(nil)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, topic := range resp.Topics {
		resources = append(resources, &SNSTopic{
			svc: n.Service,
			id:  topic.TopicArn,
		})
	}
	return resources, nil
}

type SNSTopic struct {
	svc *sns.SNS
	id  *string
}

func (topic *SNSTopic) Remove() error {
	_, err := topic.svc.DeleteTopic(&sns.DeleteTopicInput{
		TopicArn: topic.id,
	})
	return err
}

func (topic *SNSTopic) String() string {
	return fmt.Sprintf("TopicARN: %s", *topic.id)
}
