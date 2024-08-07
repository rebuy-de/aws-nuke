package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AthenaPreparedStatement struct {
	svc       *athena.Athena
	workGroup *string
	name      *string
}

func init() {
	register("AthenaPreparedStatement", ListAthenaPreparedStatements)
}

func ListAthenaPreparedStatements(sess *session.Session) ([]Resource, error) {
	svc := athena.New(sess)
	resources := []Resource{}

	workgroups, err := svc.ListWorkGroups(&athena.ListWorkGroupsInput{})
	if err != nil {
		return nil, err
	}

	for _, workgroup := range workgroups.WorkGroups {
		params := &athena.ListPreparedStatementsInput{
			WorkGroup:  workgroup.Name,
			MaxResults: aws.Int64(50),
		}

		for {
			output, err := svc.ListPreparedStatements(params)
			if err != nil {
				return nil, err
			}

			for _, statement := range output.PreparedStatements {
				resources = append(resources, &AthenaPreparedStatement{
					svc:       svc,
					workGroup: workgroup.Name,
					name:      statement.StatementName,
				})
			}

			if output.NextToken == nil {
				break
			}

			params.NextToken = output.NextToken
		}
	}

	return resources, nil
}

func (f *AthenaPreparedStatement) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("StatementName", f.name)
	properties.Set("WorkGroup", f.workGroup)

	return properties
}

func (f *AthenaPreparedStatement) Remove() error {

	_, err := f.svc.DeletePreparedStatement(&athena.DeletePreparedStatementInput{
		StatementName: f.name,
		WorkGroup:     f.workGroup,
	})

	return err
}

func (f *AthenaPreparedStatement) String() string {
	return *f.name
}
