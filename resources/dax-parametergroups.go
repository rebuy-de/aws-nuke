package resources

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dax"
)

type DAXParameterGroup struct {
	svc                *dax.DAX
	parameterGroupName *string
}

func init() {
	register("DAXParameterGroup", ListDAXParameterGroups)
}

func ListDAXParameterGroups(sess *session.Session) ([]Resource, error) {
	svc := dax.New(sess)
	resources := []Resource{}

	params := &dax.DescribeParameterGroupsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeParameterGroups(params)
		if err != nil {
			return nil, err
		}

		for _, parameterGroup := range output.ParameterGroups {
			//Ensure default is not deleted
			if !strings.Contains(*parameterGroup.ParameterGroupName, "default") {
				resources = append(resources, &DAXParameterGroup{
					svc:                svc,
					parameterGroupName: parameterGroup.ParameterGroupName,
				})
			}
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *DAXParameterGroup) Remove() error {

	_, err := f.svc.DeleteParameterGroup(&dax.DeleteParameterGroupInput{
		ParameterGroupName: f.parameterGroupName,
	})

	return err
}

func (f *DAXParameterGroup) String() string {
	return *f.parameterGroupName
}
