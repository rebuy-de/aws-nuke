package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/v2/mocks/mock_cloudformationiface"
	"github.com/stretchr/testify/assert"
)

func TestCloudformationStack_Remove_StackAlreadyDeleted(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	stack := CloudFormationStack{
		svc: mockCloudformation,
		stack: &cloudformation.Stack{
			StackName: aws.String("foobar"),
		},
	}

	mockCloudformation.EXPECT().DescribeStacks(gomock.Eq(&cloudformation.DescribeStacksInput{
		StackName: aws.String("foobar"),
	})).Return(&cloudformation.DescribeStacksOutput{
		Stacks: []*cloudformation.Stack{
			{
				StackStatus: aws.String(cloudformation.StackStatusDeleteComplete),
			},
		},
	}, nil)

	err := stack.Remove()
	a.Nil(err)
}

func TestCloudformationStack_Remove_StackDoesNotExist(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	stack := CloudFormationStack{
		svc: mockCloudformation,
		stack: &cloudformation.Stack{
			StackName: aws.String("foobar"),
		},
	}

	mockCloudformation.EXPECT().DescribeStacks(gomock.Eq(&cloudformation.DescribeStacksInput{
		StackName: aws.String("foobar"),
	})).Return(nil, awserr.New("ValidationFailed", "Stack with id foobar does not exist", nil))

	err := stack.Remove()
	a.Nil(err)
}

func TestCloudformationStack_Remove_DeleteFailed(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	stack := CloudFormationStack{
		svc: mockCloudformation,
		stack: &cloudformation.Stack{
			StackName: aws.String("foobar"),
		},
	}

	gomock.InOrder(
		mockCloudformation.EXPECT().DescribeStacks(gomock.Eq(&cloudformation.DescribeStacksInput{
			StackName: aws.String("foobar"),
		})).Return(&cloudformation.DescribeStacksOutput{
			Stacks: []*cloudformation.Stack{
				{
					StackStatus: aws.String(cloudformation.StackStatusDeleteFailed),
				},
			},
		}, nil),
		mockCloudformation.EXPECT().ListStackResources(gomock.Eq(&cloudformation.ListStackResourcesInput{
			StackName: aws.String("foobar"),
		})).Return(&cloudformation.ListStackResourcesOutput{
			StackResourceSummaries: []*cloudformation.StackResourceSummary{
				{
					ResourceStatus:    aws.String(cloudformation.ResourceStatusDeleteComplete),
					LogicalResourceId: aws.String("fooDeleteComplete"),
				},
				{
					ResourceStatus:    aws.String(cloudformation.ResourceStatusDeleteFailed),
					LogicalResourceId: aws.String("fooDeleteFailed"),
				},
			},
		}, nil),
		mockCloudformation.EXPECT().DeleteStack(gomock.Eq(&cloudformation.DeleteStackInput{
			StackName: aws.String("foobar"),
			RetainResources: []*string{
				aws.String("fooDeleteFailed"),
			},
		})).Return(nil, nil),
		mockCloudformation.EXPECT().WaitUntilStackDeleteComplete(gomock.Eq(&cloudformation.DescribeStacksInput{
			StackName: aws.String("foobar"),
		})).Return(nil),
	)

	err := stack.Remove()
	a.Nil(err)
}

// if the stack is currently in delete in progress
func TestCloudformationStack_Remove_DeleteInProgress(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	stack := CloudFormationStack{
		svc: mockCloudformation,
		stack: &cloudformation.Stack{
			StackName: aws.String("foobar"),
		},
	}

	gomock.InOrder(
		mockCloudformation.EXPECT().DescribeStacks(gomock.Eq(&cloudformation.DescribeStacksInput{
			StackName: aws.String("foobar"),
		})).Return(&cloudformation.DescribeStacksOutput{
			Stacks: []*cloudformation.Stack{
				{
					StackStatus: aws.String(cloudformation.StackStatusDeleteInProgress),
				},
			},
		}, nil),

		mockCloudformation.EXPECT().WaitUntilStackDeleteComplete(gomock.Eq(&cloudformation.DescribeStacksInput{
			StackName: aws.String("foobar"),
		})).Return(nil),
	)

	err := stack.Remove()
	a.Nil(err)
}

func TestCloudformationStack_Remove_Stack_InCompletedStatus(t *testing.T) {
	tests := []string{
		cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusReviewInProgress,
		cloudformation.StackStatusRollbackComplete,
		cloudformation.StackStatusRollbackFailed,
		cloudformation.StackStatusUpdateComplete,
		cloudformation.StackStatusUpdateRollbackComplete,
		cloudformation.StackStatusUpdateRollbackFailed,
	}

	for _, stackStatus := range tests {
		t.Run(stackStatus, func(t *testing.T) {
			a := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

			stack := CloudFormationStack{
				svc: mockCloudformation,
				stack: &cloudformation.Stack{
					StackName: aws.String("foobar"),
				},
			}

			gomock.InOrder(
				mockCloudformation.EXPECT().DescribeStacks(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(&cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackStatus: aws.String(stackStatus),
						},
					},
				}, nil),

				mockCloudformation.EXPECT().DeleteStack(gomock.Eq(&cloudformation.DeleteStackInput{
					StackName: aws.String("foobar"),
				})).Return(nil, nil),

				mockCloudformation.EXPECT().WaitUntilStackDeleteComplete(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(nil),
			)

			err := stack.Remove()
			a.Nil(err)
		})
	}
}

func TestCloudformationStack_Remove_Stack_CreateInProgress(t *testing.T) {
	tests := []string{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusRollbackInProgress,
	}

	for _, stackStatus := range tests {
		t.Run(stackStatus, func(t *testing.T) {
			a := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

			stack := CloudFormationStack{
				svc: mockCloudformation,
				stack: &cloudformation.Stack{
					StackName: aws.String("foobar"),
				},
			}

			gomock.InOrder(
				mockCloudformation.EXPECT().DescribeStacks(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(&cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackStatus: aws.String(stackStatus),
						},
					},
				}, nil),

				mockCloudformation.EXPECT().WaitUntilStackCreateComplete(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(nil),

				mockCloudformation.EXPECT().DeleteStack(gomock.Eq(&cloudformation.DeleteStackInput{
					StackName: aws.String("foobar"),
				})).Return(nil, nil),

				mockCloudformation.EXPECT().WaitUntilStackDeleteComplete(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(nil),
			)

			err := stack.Remove()
			a.Nil(err)
		})
	}
}

func TestCloudformationStack_Remove_Stack_UpdateInProgress(t *testing.T) {
	tests := []string{
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackInProgress,
	}

	for _, stackStatus := range tests {
		t.Run(stackStatus, func(t *testing.T) {
			a := assert.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

			stack := CloudFormationStack{
				svc: mockCloudformation,
				stack: &cloudformation.Stack{
					StackName: aws.String("foobar"),
				},
			}

			gomock.InOrder(
				mockCloudformation.EXPECT().DescribeStacks(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(&cloudformation.DescribeStacksOutput{
					Stacks: []*cloudformation.Stack{
						{
							StackStatus: aws.String(stackStatus),
						},
					},
				}, nil),

				mockCloudformation.EXPECT().WaitUntilStackUpdateComplete(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(nil),

				mockCloudformation.EXPECT().DeleteStack(gomock.Eq(&cloudformation.DeleteStackInput{
					StackName: aws.String("foobar"),
				})).Return(nil, nil),

				mockCloudformation.EXPECT().WaitUntilStackDeleteComplete(gomock.Eq(&cloudformation.DescribeStacksInput{
					StackName: aws.String("foobar"),
				})).Return(nil),
			)

			err := stack.Remove()
			a.Nil(err)
		})
	}
}
