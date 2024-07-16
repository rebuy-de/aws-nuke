package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type GlueSession struct {
	svc *glue.Glue
	id  *string
}

func init() {
	register("GlueSession", ListGlueSessions)
}

func ListGlueSessions(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.ListSessionsInput{
		MaxResults: aws.Int64(25),
	}

	for {
		output, err := svc.ListSessions(params)
		if err != nil {
			return nil, err
		}

		for _, session := range output.Sessions {
			resources = append(resources, &GlueSession{
				svc: svc,
				id:  session.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueSession) Remove() error {
	_, err := f.svc.DeleteSession(&glue.DeleteSessionInput{
		Id: f.id,
	})

	return err
}

func (f *GlueSession) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Id", f.id)

	return properties
}

func (f *GlueSession) String() string {
	return *f.id
}
