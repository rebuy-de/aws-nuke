package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type LexBot struct {
	svc    *lexmodelbuildingservice.LexModelBuildingService
	name   *string
	status *string
}

func init() {
	register("LexBot", ListLexBots)
}

func ListLexBots(sess *session.Session) ([]Resource, error) {
	svc := lexmodelbuildingservice.New(sess)
	resources := []Resource{}

	params := &lexmodelbuildingservice.GetBotsInput{
		MaxResults: aws.Int64(10),
	}

	for {
		output, err := svc.GetBots(params)
		if err != nil {
			return nil, err
		}

		for _, bot := range output.Bots {
			resources = append(resources, &LexBot{
				svc:    svc,
				name:   bot.Name,
				status: bot.Status,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *LexBot) Remove() error {

	_, err := f.svc.DeleteBot(&lexmodelbuildingservice.DeleteBotInput{
		Name: f.name,
	})

	return err
}

func (f *LexBot) String() string {
	return *f.name
}

func (f *LexBot) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("Name", f.name).
		Set("Status", f.status)
	return properties
}
