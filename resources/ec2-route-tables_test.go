package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/mocks/mock_ec2iface"
	"github.com/stretchr/testify/assert"
)

func TestEC2RouteTable_Delete(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2 := mock_ec2iface.NewMockEC2API(ctrl)

	routeTable := EC2RouteTable{
		svc: mockEc2,
		routeTable: &ec2.RouteTable{
			RouteTableId: aws.String("foo1"),
		},
	}

	mockEc2.EXPECT().DeleteRouteTable(gomock.Eq(&ec2.DeleteRouteTableInput{
		RouteTableId: aws.String("foo1"),
	})).Return(&ec2.DeleteRouteTableOutput{}, nil)

	err := routeTable.Remove()
	a.Nil(err)
}

func TestEC2RouteTable_Delete_NotFound(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2 := mock_ec2iface.NewMockEC2API(ctrl)

	routeTable := EC2RouteTable{
		svc: mockEc2,
		routeTable: &ec2.RouteTable{
			RouteTableId: aws.String("foo1"),
		},
	}

	mockEc2.EXPECT().DeleteRouteTable(gomock.Eq(&ec2.DeleteRouteTableInput{
		RouteTableId: aws.String("foo1"),
	})).Return(nil, awserr.New("InvalidRouteTableID.NotFound", "foo", nil))

	err := routeTable.Remove()
	a.Nil(err)
}
