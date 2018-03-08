package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudhsmv2"
)

type CloudHSMV2ClusterHSM struct {
	svc       *cloudhsmv2.CloudHSMV2
	clusterID *string
	hsmID     *string
}

func init() {
	register("CloudHSMV2ClusterHSM", ListCloudHSMV2ClusterHSMs)
}

func ListCloudHSMV2ClusterHSMs(sess *session.Session) ([]Resource, error) {
	svc := cloudhsmv2.New(sess)
	resources := []Resource{}

	params := &cloudhsmv2.DescribeClustersInput{
		MaxResults: aws.Int64(25),
	}

	for {
		resp, err := svc.DescribeClusters(params)
		if err != nil {
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			for _, hsm := range cluster.Hsms {
				resources = append(resources, &CloudHSMV2ClusterHSM{
					svc:       svc,
					clusterID: hsm.ClusterId,
					hsmID:     hsm.HsmId,
				})
			}

		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CloudHSMV2ClusterHSM) Remove() error {

	_, err := f.svc.DeleteHsm(&cloudhsmv2.DeleteHsmInput{
		ClusterId: f.clusterID,
		HsmId:     f.hsmID,
	})

	return err
}

func (f *CloudHSMV2ClusterHSM) String() string {
	return *f.hsmID
}
