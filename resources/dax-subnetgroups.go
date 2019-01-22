package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dax"
	"fmt"
)

type DAXSubnetGroup struct {
	svc             *dax.DAX
	subnetGroupName *string
}

func init() {
	register("DAXSubnetGroup", ListDAXSubnetGroups)
}

func ListDAXSubnetGroups(sess *session.Session) ([]Resource, error) {
	svc := dax.New(sess)
	resources := []Resource{}

	params := &dax.DescribeSubnetGroupsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeSubnetGroups(params)
		if err != nil {
			return nil, err
		}

		for _, subnet := range output.SubnetGroups {
			resources = append(resources, &DAXSubnetGroup{
				svc:             svc,
				subnetGroupName: subnet.SubnetGroupName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *DAXSubnetGroup) Filter() error {
	if *f.subnetGroupName == "default" {
		return fmt.Errorf("Cannot delete default DAX Subnet group")
	}
	return nil
}

func (f *DAXSubnetGroup) Remove() error {

	_, err := f.svc.DeleteSubnetGroup(&dax.DeleteSubnetGroupInput{
		SubnetGroupName: f.subnetGroupName,
	})

	return err
}

func (f *DAXSubnetGroup) String() string {
	return *f.subnetGroupName
}
