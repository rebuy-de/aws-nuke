package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type LexIntent struct {
	svc  *lexmodelbuildingservice.LexModelBuildingService
	name *string
}

func init() {
	register("LexIntent", ListLexIntents)
}

func ListLexIntents(sess *session.Session) ([]Resource, error) {
	svc := lexmodelbuildingservice.New(sess)
	resources := []Resource{}

	params := &lexmodelbuildingservice.GetIntentsInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.GetIntents(params)
		if err != nil {
			return nil, err
		}

		for _, bot := range output.Intents {
			resources = append(resources, &LexIntent{
				svc:  svc,
				name: bot.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *LexIntent) Remove() error {

	_, err := f.svc.DeleteIntent(&lexmodelbuildingservice.DeleteIntentInput{
		Name: f.name,
	})

	return err
}

func (f *LexIntent) String() string {
	return *f.name
}

func (f *LexIntent) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("Name", f.name)
	return properties
}
