package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RedshiftServerlessWorkgroup struct {
	svc       *redshiftserverless.RedshiftServerless
	workgroup *redshiftserverless.Workgroup
}

func init() {
	register("RedshiftServerlessWorkgroup", ListRedshiftServerlessWorkgroups)
}

func ListRedshiftServerlessWorkgroups(sess *session.Session) ([]Resource, error) {
	svc := redshiftserverless.New(sess)
	resources := []Resource{}

	params := &redshiftserverless.ListWorkgroupsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListWorkgroups(params)
		if err != nil {
			return nil, err
		}

		for _, workgroup := range output.Workgroups {
			resources = append(resources, &RedshiftServerlessWorkgroup{
				svc:       svc,
				workgroup: workgroup,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (w *RedshiftServerlessWorkgroup) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreationDate", w.workgroup.CreationDate).
		Set("Namespace", w.workgroup.NamespaceName).
		Set("WorkgroupName", w.workgroup.WorkgroupName)

	return properties
}

func (w *RedshiftServerlessWorkgroup) Remove() error {
	_, err := w.svc.DeleteWorkgroup(&redshiftserverless.DeleteWorkgroupInput{
		WorkgroupName: w.workgroup.WorkgroupName,
	})

	return err
}

func (w *RedshiftServerlessWorkgroup) String() string {
	return *w.workgroup.WorkgroupName
}
