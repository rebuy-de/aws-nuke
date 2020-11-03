package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/mocks/mock_ec2iface"
	"github.com/stretchr/testify/assert"
)

func TestEC2EgressOnlyInternetGateway_Remove_Successful(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEC2 := mock_ec2iface.NewMockEC2API(ctrl)

	egressOnlyIGW := EC2EgressOnlyInternetGateway{
		svc: mockEC2,
		igw: &ec2.EgressOnlyInternetGateway{
			EgressOnlyInternetGatewayId: aws.String("foobar"),
		},
	}

	mockEC2.EXPECT().DeleteEgressOnlyInternetGateway(gomock.Eq(&ec2.DeleteEgressOnlyInternetGatewayInput{
		EgressOnlyInternetGatewayId: aws.String("foobar"),
	})).Return(&ec2.DeleteEgressOnlyInternetGatewayOutput{
		ReturnCode: aws.Bool(true),
	}, nil)

	err := egressOnlyIGW.Remove()
	a.Nil(err)
}
