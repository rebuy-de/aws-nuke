package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type SNSTopic struct {
	svc  *sns.SNS
	id   *string
	tags []*sns.Tag
}

func init() {
	register("SNSTopic", ListSNSTopics)
}

func ListSNSTopics(sess *session.Session) ([]Resource, error) {
	svc := sns.New(sess)

	topics := make([]*sns.Topic, 0)

	params := &sns.ListTopicsInput{}

	err := svc.ListTopicsPages(params, func(page *sns.ListTopicsOutput, lastPage bool) bool {
		for _, out := range page.Topics {
			topics = append(topics, out)
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, topic := range topics {
		tags, err := svc.ListTagsForResource(&sns.ListTagsForResourceInput{
			ResourceArn: topic.TopicArn,
		})

		if err != nil {
			continue
		}

		resources = append(resources, &SNSTopic{
			svc:  svc,
			id:   topic.TopicArn,
			tags: tags.Tags,
		})
	}
	return resources, nil
}

func (topic *SNSTopic) Remove() error {
	_, err := topic.svc.DeleteTopic(&sns.DeleteTopicInput{
		TopicArn: topic.id,
	})
	return err
}

func (topic *SNSTopic) Properties() types.Properties {
	properties := types.NewProperties()

	for _, tag := range topic.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.Set("TopicARN", topic.id)

	return properties
}

func (topic *SNSTopic) String() string {
	return fmt.Sprintf("TopicARN: %s", *topic.id)
}
