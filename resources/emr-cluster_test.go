package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/mocks/mock_emriface"
	"github.com/rebuy-de/aws-nuke/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestEMRCluster_Remove_TerminationProtectionEnabled(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEMRCluster := mock_emriface.NewMockEMRAPI(ctrl)

	emrCluster := EMRCluster{
		svc:   mockEMRCluster,
		ID:    aws.String("foobar"),
		state: aws.String(emr.ClusterStateRunning),
		featureFlags: config.FeatureFlags{
			DisableDeletionProtection: config.DisableDeletionProtection{
				EMRCluster: true,
			},
		},
	}

	gomock.InOrder(
		mockEMRCluster.EXPECT().TerminateJobFlows(gomock.Eq(&emr.TerminateJobFlowsInput{
			JobFlowIds: aws.StringSlice([]string{"foobar"}),
		})).Return(nil, awserr.New("ValidationException", "Could not shut down one or more job flows since they are termination protected.", nil)),
		mockEMRCluster.EXPECT().SetTerminationProtection(gomock.Eq(&emr.SetTerminationProtectionInput{
			JobFlowIds:           aws.StringSlice([]string{"foobar"}),
			TerminationProtected: aws.Bool(false),
		})),
		mockEMRCluster.EXPECT().TerminateJobFlows(gomock.Eq(&emr.TerminateJobFlowsInput{
			JobFlowIds: aws.StringSlice([]string{"foobar"}),
		})).Return(nil, nil),
	)

	err := emrCluster.Remove()
	a.Nil(err)
}

func TestEMRCluster_Remove_TerminationProtectionDisabled(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEMRCluster := mock_emriface.NewMockEMRAPI(ctrl)

	emrCluster := EMRCluster{
		svc:   mockEMRCluster,
		ID:    aws.String("foobar"),
		state: aws.String(emr.ClusterStateRunning),
		featureFlags: config.FeatureFlags{
			DisableDeletionProtection: config.DisableDeletionProtection{
				EMRCluster: false,
			},
		},
	}

	mockEMRCluster.EXPECT().TerminateJobFlows(gomock.Eq(&emr.TerminateJobFlowsInput{
		JobFlowIds: aws.StringSlice([]string{"foobar"}),
	})).Return(nil, nil)

	err := emrCluster.Remove()
	a.Nil(err)
}

func TestEMRCluster_DisableProtection(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEMRCluster := mock_emriface.NewMockEMRAPI(ctrl)

	emrCluster := EMRCluster{
		svc: mockEMRCluster,
		ID:  aws.String("foobar"),
	}

	mockEMRCluster.EXPECT().SetTerminationProtection(gomock.Eq(&emr.SetTerminationProtectionInput{
		JobFlowIds:           aws.StringSlice([]string{"foobar"}),
		TerminationProtected: aws.Bool(false),
	})).Return(nil, nil)

	err := emrCluster.DisableProtection()
	a.Nil(err)
}
