package resources

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCloudControlParseProperties(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	cases := []struct {
		name    string
		payload string
		want    string
	}{
		{
			name:    "ActualEC2VPC",
			payload: `{"VpcId":"vpc-456","InstanceTenancy":"default","CidrBlockAssociations":["vpc-cidr-assoc-1234", "vpc-cidr-assoc-5678"],"CidrBlock":"10.10.0.0/16","Tags":[{"Value":"Kubernetes VPC","Key":"Name"}]}`,
			want:    `[CidrBlock: "10.10.0.0/16", CidrBlockAssociations.["vpc-cidr-assoc-1234"]: "true", CidrBlockAssociations.["vpc-cidr-assoc-5678"]: "true", InstanceTenancy: "default", Tags.["Name"]: "Kubernetes VPC", VpcId: "vpc-456"]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := cloudControlParseProperties(tc.payload)
			require.NoError(t, err)
			require.Equal(t, tc.want, result.String())
		})
	}
}
