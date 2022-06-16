package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/v2/mocks/mock_cloudformationiface"
	"github.com/stretchr/testify/assert"
)

func TestCloudformationStackSet_Remove(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	stackSet := CloudFormationStackSet{
		svc: mockCloudformation,
		stackSetSummary: &cloudformation.StackSetSummary{
			StackSetName: aws.String("foobar"),
		},
	}

	mockCloudformation.EXPECT().ListStackInstances(gomock.Eq(&cloudformation.ListStackInstancesInput{
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.ListStackInstancesOutput{
		Summaries: []*cloudformation.StackInstanceSummary{
			{
				Account: aws.String("a1"),
				Region:  aws.String("r1"),
			},
			{
				Account: aws.String("a1"),
				Region:  aws.String("r2"),
			},
		},
	}, nil)

	mockCloudformation.EXPECT().DeleteStackInstances(gomock.Eq(&cloudformation.DeleteStackInstancesInput{
		StackSetName: aws.String("foobar"),
		Accounts:     []*string{aws.String("a1")},
		Regions:      []*string{aws.String("r1"), aws.String("r2")},
		RetainStacks: aws.Bool(true),
	})).Return(&cloudformation.DeleteStackInstancesOutput{
		OperationId: aws.String("o1"),
	}, nil)

	describeStackSetStatuses := []string{
		cloudformation.StackSetOperationResultStatusPending,
		cloudformation.StackSetOperationResultStatusRunning,
		cloudformation.StackSetOperationResultStatusSucceeded,
	}
	describeStackSetOperationCalls := make([]*gomock.Call, len(describeStackSetStatuses))
	for i, status := range describeStackSetStatuses {
		describeStackSetOperationCalls[i] = mockCloudformation.EXPECT().DescribeStackSetOperation(gomock.Eq(&cloudformation.DescribeStackSetOperationInput{
			OperationId:  aws.String("o1"),
			StackSetName: aws.String("foobar"),
		})).Return(&cloudformation.DescribeStackSetOperationOutput{
			StackSetOperation: &cloudformation.StackSetOperation{
				Status: aws.String(status),
			},
		}, nil)
	}
	gomock.InOrder(describeStackSetOperationCalls...)

	mockCloudformation.EXPECT().DeleteStackSet(gomock.Eq(&cloudformation.DeleteStackSetInput{
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.DeleteStackSetOutput{}, nil)

	err := stackSet.Remove()
	a.Nil(err)
}

func TestCloudformationStackSet_Remove_MultipleAccounts(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	stackSet := CloudFormationStackSet{
		svc: mockCloudformation,
		stackSetSummary: &cloudformation.StackSetSummary{
			StackSetName: aws.String("foobar"),
		},
	}

	mockCloudformation.EXPECT().ListStackInstances(gomock.Eq(&cloudformation.ListStackInstancesInput{
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.ListStackInstancesOutput{
		Summaries: []*cloudformation.StackInstanceSummary{
			{
				Account: aws.String("a1"),
				Region:  aws.String("r1"),
			},
			{
				Account: aws.String("a1"),
				Region:  aws.String("r2"),
			},
			{
				Account: aws.String("a2"),
				Region:  aws.String("r2"),
			},
		},
	}, nil)

	mockCloudformation.EXPECT().DeleteStackInstances(gomock.Eq(&cloudformation.DeleteStackInstancesInput{
		StackSetName: aws.String("foobar"),
		Accounts:     []*string{aws.String("a1")},
		Regions:      []*string{aws.String("r1"), aws.String("r2")},
		RetainStacks: aws.Bool(true),
	})).Return(&cloudformation.DeleteStackInstancesOutput{
		OperationId: aws.String("a1-oId"),
	}, nil)
	mockCloudformation.EXPECT().DeleteStackInstances(gomock.Eq(&cloudformation.DeleteStackInstancesInput{
		StackSetName: aws.String("foobar"),
		Accounts:     []*string{aws.String("a2")},
		Regions:      []*string{aws.String("r2")},
		RetainStacks: aws.Bool(true),
	})).Return(&cloudformation.DeleteStackInstancesOutput{
		OperationId: aws.String("a2-oId"),
	}, nil)

	mockCloudformation.EXPECT().DescribeStackSetOperation(gomock.Eq(&cloudformation.DescribeStackSetOperationInput{
		OperationId:  aws.String("a1-oId"),
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.DescribeStackSetOperationOutput{
		StackSetOperation: &cloudformation.StackSetOperation{
			Status: aws.String(cloudformation.StackSetOperationResultStatusSucceeded),
		},
	}, nil)
	mockCloudformation.EXPECT().DescribeStackSetOperation(gomock.Eq(&cloudformation.DescribeStackSetOperationInput{
		OperationId:  aws.String("a2-oId"),
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.DescribeStackSetOperationOutput{
		StackSetOperation: &cloudformation.StackSetOperation{
			Status: aws.String(cloudformation.StackSetOperationResultStatusSucceeded),
		},
	}, nil)

	mockCloudformation.EXPECT().DeleteStackSet(gomock.Eq(&cloudformation.DeleteStackSetInput{
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.DeleteStackSetOutput{}, nil)

	err := stackSet.Remove()
	a.Nil(err)
}

func TestCloudformationStackSet_Remove_DeleteStackInstanceFailed(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	stackSet := CloudFormationStackSet{
		svc: mockCloudformation,
		stackSetSummary: &cloudformation.StackSetSummary{
			StackSetName: aws.String("foobar"),
		},
	}

	mockCloudformation.EXPECT().ListStackInstances(gomock.Eq(&cloudformation.ListStackInstancesInput{
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.ListStackInstancesOutput{
		Summaries: []*cloudformation.StackInstanceSummary{
			{
				Account: aws.String("a1"),
				Region:  aws.String("r1"),
			},
		},
	}, nil)

	mockCloudformation.EXPECT().DeleteStackInstances(gomock.Eq(&cloudformation.DeleteStackInstancesInput{
		StackSetName: aws.String("foobar"),
		Accounts:     []*string{aws.String("a1")},
		Regions:      []*string{aws.String("r1")},
		RetainStacks: aws.Bool(true),
	})).Return(&cloudformation.DeleteStackInstancesOutput{
		OperationId: aws.String("o1"),
	}, nil)

	mockCloudformation.EXPECT().DescribeStackSetOperation(gomock.Eq(&cloudformation.DescribeStackSetOperationInput{
		OperationId:  aws.String("o1"),
		StackSetName: aws.String("foobar"),
	})).Return(&cloudformation.DescribeStackSetOperationOutput{
		StackSetOperation: &cloudformation.StackSetOperation{
			Status: aws.String(cloudformation.StackSetOperationResultStatusFailed),
		},
	}, nil)

	err := stackSet.Remove()
	a.EqualError(err, "unable to delete stackSet=foobar operationId=o1 status=FAILED")
}
