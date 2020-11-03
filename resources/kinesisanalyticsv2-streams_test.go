package resources

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/mocks/mock_kinesisanalyticsv2iface"
	"github.com/stretchr/testify/assert"
)

func TestKinesisAnalyticsApplicationV2_Remove(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKinesisAnalyticsV2 := mock_kinesisanalyticsv2iface.NewMockKinesisAnalyticsV2API(ctrl)

	kinesisAnalyticsApplicationV2 := KinesisAnalyticsApplicationV2{
		svc:             mockKinesisAnalyticsV2,
		applicationName: aws.String("foobar"),
	}

	now := time.Now()

	gomock.InOrder(
		mockKinesisAnalyticsV2.EXPECT().DescribeApplication(gomock.Eq(&kinesisanalyticsv2.DescribeApplicationInput{
			ApplicationName: aws.String("foobar"),
		})).Return(&kinesisanalyticsv2.DescribeApplicationOutput{
			ApplicationDetail: &kinesisanalyticsv2.ApplicationDetail{
				CreateTimestamp: aws.Time(now),
			},
		}, nil),
		mockKinesisAnalyticsV2.EXPECT().DeleteApplication(gomock.Eq(&kinesisanalyticsv2.DeleteApplicationInput{
			ApplicationName: aws.String("foobar"),
			CreateTimestamp: aws.Time(now),
		})).Return(nil, nil),
	)

	err := kinesisAnalyticsApplicationV2.Remove()
	a.Nil(err)
}
