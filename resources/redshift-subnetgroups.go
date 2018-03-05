package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
)

type RedshiftSubnetGroup struct {
	svc                    *redshift.Redshift
	clusterSubnetGroupName *string
}

func init() {
	register("RedshiftSubnetGroup", ListRedshiftSubnetGroups)
}

func ListRedshiftSubnetGroups(sess *session.Session) ([]Resource, error) {
	svc := redshift.New(sess)
	resources := []Resource{}

	params := &redshift.DescribeClusterSubnetGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeClusterSubnetGroups(params)
		if err != nil {
			return nil, err
		}

		for _, subnetGroup := range output.ClusterSubnetGroups {
			resources = append(resources, &RedshiftSubnetGroup{
				svc: svc,
				clusterSubnetGroupName: subnetGroup.ClusterSubnetGroupName,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *RedshiftSubnetGroup) Remove() error {

	_, err := f.svc.DeleteClusterSubnetGroup(&redshift.DeleteClusterSubnetGroupInput{
		ClusterSubnetGroupName: f.clusterSubnetGroupName,
	})

	return err
}

func (f *RedshiftSubnetGroup) String() string {
	return *f.clusterSubnetGroupName
}
