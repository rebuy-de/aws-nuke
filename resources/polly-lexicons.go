package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type PollyLexicon struct {
	svc        *polly.Polly
	name       *string
	attributes *polly.LexiconAttributes
}

func init() {
	register("PollyLexicon", ListPollyLexicons)
}

func ListPollyLexicons(sess *session.Session) ([]Resource, error) {
	svc := polly.New(sess)
	resources := []Resource{}

	listLexiconsInput := &polly.ListLexiconsInput{}

	listOutput, err := svc.ListLexicons(listLexiconsInput)
	if err != nil {
		return nil, err
	}
	for _, lexicon := range listOutput.Lexicons {
		resources = append(resources, &PollyLexicon{
			svc:        svc,
			name:       lexicon.Name,
			attributes: lexicon.Attributes,
		})
	}
	return resources, nil
}

func (lexicon *PollyLexicon) Remove() error {
	deleteInput := &polly.DeleteLexiconInput{
		Name: lexicon.name,
	}
	_, err := lexicon.svc.DeleteLexicon(deleteInput)
	return err
}

func (lexicon *PollyLexicon) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", lexicon.name)
	properties.Set("Alphabet", lexicon.attributes.Alphabet)
	properties.Set("LanguageCode", lexicon.attributes.LanguageCode)
	properties.Set("LastModified", lexicon.attributes.LastModified)
	properties.Set("LexemesCount", lexicon.attributes.LexemesCount)
	properties.Set("LexiconArn", lexicon.attributes.LexiconArn)
	properties.Set("Size", lexicon.attributes.Size)
	return properties
}
