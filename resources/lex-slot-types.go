package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type LexSlotType struct {
	svc  *lexmodelbuildingservice.LexModelBuildingService
	name *string
}

func init() {
	register("LexSlotType", ListLexSlotTypes)
}

func ListLexSlotTypes(sess *session.Session) ([]Resource, error) {
	svc := lexmodelbuildingservice.New(sess)
	resources := []Resource{}

	params := &lexmodelbuildingservice.GetSlotTypesInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.GetSlotTypes(params)
		if err != nil {
			return nil, err
		}

		for _, bot := range output.SlotTypes {
			resources = append(resources, &LexSlotType{
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

func (f *LexSlotType) Remove() error {

	_, err := f.svc.DeleteSlotType(&lexmodelbuildingservice.DeleteSlotTypeInput{
		Name: f.name,
	})

	return err
}

func (f *LexSlotType) String() string {
	return *f.name
}

func (f *LexSlotType) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("Name", f.name)
	return properties
}
