package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
)

type RedshiftParameterGroup struct {
	svc                *redshift.Redshift
	parameterGroupName *string
}

func init() {
	register("RedshiftParameterGroup", ListRedshiftParameterGroup)
}

func ListRedshiftParameterGroup(sess *session.Session) ([]Resource, error) {
	svc := redshift.New(sess)
	resources := []Resource{}

	params := &redshift.DescribeClusterParameterGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeClusterParameterGroups(params)
		if err != nil {
			return nil, err
		}

		for _, parameterGroup := range output.ParameterGroups {
			if !strings.Contains(*parameterGroup.ParameterGroupName, "default.redshift") {
				resources = append(resources, &RedshiftParameterGroup{
					svc:                svc,
					parameterGroupName: parameterGroup.ParameterGroupName,
				})
			}
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *RedshiftParameterGroup) Remove() error {

	_, err := f.svc.DeleteClusterParameterGroup(&redshift.DeleteClusterParameterGroupInput{
		ParameterGroupName: f.parameterGroupName,
	})

	return err
}

func (f *RedshiftParameterGroup) String() string {
	return *f.parameterGroupName
}
