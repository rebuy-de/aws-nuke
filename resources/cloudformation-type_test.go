package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/v2/mocks/mock_cloudformationiface"
	"github.com/stretchr/testify/assert"
)

func TestCloudformationType_Remove(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloudformation := mock_cloudformationiface.NewMockCloudFormationAPI(ctrl)

	cfnType := CloudFormationType{
		svc: mockCloudformation,
		typeSummary: &cloudformation.TypeSummary{
			TypeArn: aws.String("foobar"),
		},
	}

	mockCloudformation.EXPECT().ListTypeVersions(gomock.Eq(&cloudformation.ListTypeVersionsInput{
		Arn: aws.String("foobar"),
	})).Return(&cloudformation.ListTypeVersionsOutput{
		TypeVersionSummaries: []*cloudformation.TypeVersionSummary{
			{
				IsDefaultVersion: aws.Bool(true),
				VersionId:        aws.String("1"),
				TypeName:         aws.String("t1"),
				Type:             aws.String("RESOURCE"),
			},
		},
		NextToken: aws.String("nextToken"),
	}, nil)
	mockCloudformation.EXPECT().ListTypeVersions(gomock.Eq(&cloudformation.ListTypeVersionsInput{
		Arn:       aws.String("foobar"),
		NextToken: aws.String("nextToken"),
	})).Return(&cloudformation.ListTypeVersionsOutput{
		TypeVersionSummaries: []*cloudformation.TypeVersionSummary{
			{
				IsDefaultVersion: aws.Bool(false),
				VersionId:        aws.String("2"),
				TypeName:         aws.String("t2"),
				Type:             aws.String("RESOURCE"),
			},
		},
	}, nil)

	mockCloudformation.EXPECT().DeregisterType(&cloudformation.DeregisterTypeInput{
		VersionId: aws.String("2"),
		TypeName:  aws.String("t2"),
		Type:      aws.String("RESOURCE"),
	})
	mockCloudformation.EXPECT().DeregisterType(&cloudformation.DeregisterTypeInput{
		Arn: aws.String("foobar"),
	})

	err := cfnType.Remove()
	a.Nil(err)

}
