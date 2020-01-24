package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/athena/athenaiface"
)

type AthenaNamedQuery struct {
	svc athenaiface.AthenaAPI
	id  *string
}

func init() {
	register("AthenaNamedQuery", ListAthenaNamedQuery)
}

func ListAthenaNamedQuery(sess *session.Session) ([]Resource, error) {
	svc := athena.New(sess)
	resources := []Resource{}

	params := &athena.ListNamedQueriesInput{
		MaxResults: aws.Int64(20),
	}

	for {
		resp, err := svc.ListNamedQueries(params)
		if err != nil {
			return nil, err
		}

		for _, namedQueryId := range resp.NamedQueryIds {
			resources = append(resources, &AthenaNamedQuery{
				svc: svc,
				id:  namedQueryId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *AthenaNamedQuery) Remove() error {

	_, err := f.svc.DeleteNamedQuery(&athena.DeleteNamedQueryInput{
		NamedQueryId: f.id,
	})

	return err
}

func (f *AthenaNamedQuery) String() string {
	return *f.id
}
