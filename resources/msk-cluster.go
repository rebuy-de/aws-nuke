package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MSKCluster struct {
	svc  *kafka.Kafka
	arn  string
	name string
}

func init() {
	register("MSKCluster", ListMSKCluster)
}

func ListMSKCluster(sess *session.Session) ([]Resource, error) {
	svc := kafka.New(sess)
	params := &kafka.ListClustersInput{}
	resp, err := svc.ListClusters(params)

	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, cluster := range resp.ClusterInfoList {
		resources = append(resources, &MSKCluster{
			svc:  svc,
			arn:  *cluster.ClusterArn,
			name: *cluster.ClusterName,
		})

	}
	return resources, nil
}

func (m *MSKCluster) Remove() error {
	params := &kafka.DeleteClusterInput{
		ClusterArn: &m.arn,
	}

	_, err := m.svc.DeleteCluster(params)
	if err != nil {
		return err
	}
	return nil
}

func (m *MSKCluster) String() string {
	return m.arn
}

func (m *MSKCluster) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ARN", m.arn)
	properties.Set("Name", m.name)

	return properties
}
