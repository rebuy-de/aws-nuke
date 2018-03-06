package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mq"
)

type MQBroker struct {
	svc      *mq.MQ
	brokerID *string
}

func init() {
	register("MQBroker", ListMQBrokers)
}

func ListMQBrokers(sess *session.Session) ([]Resource, error) {
	svc := mq.New(sess)
	resources := []Resource{}

	params := &mq.ListBrokersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.ListBrokers(params)
		if err != nil {
			return nil, err
		}

		for _, broker := range resp.BrokerSummaries {
			resources = append(resources, &MQBroker{
				svc:      svc,
				brokerID: broker.BrokerId,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *MQBroker) Remove() error {

	_, err := f.svc.DeleteBroker(&mq.DeleteBrokerInput{
		BrokerId: f.brokerID,
	})

	return err
}

func (f *MQBroker) String() string {
	return *f.brokerID
}
