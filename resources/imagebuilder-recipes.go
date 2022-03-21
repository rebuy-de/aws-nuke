package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ImageBuilderRecipe struct {
	svc *imagebuilder.Imagebuilder
	arn string
}

func init() {
	register("ImageBuilderRecipe", ListImageBuilderRecipes)
}

func ListImageBuilderRecipes(sess *session.Session) ([]Resource, error) {
	svc := imagebuilder.New(sess)
	params := &imagebuilder.ListImageRecipesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListImageRecipes(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.ImageRecipeSummaryList {
			resources = append(resources, &ImageBuilderRecipe{
				svc: svc,
				arn: *out.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListImageRecipesInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func (e *ImageBuilderRecipe) Remove() error {
	_, err := e.svc.DeleteImageRecipe(&imagebuilder.DeleteImageRecipeInput{
		ImageRecipeArn: &e.arn,
	})
	return err
}

func (e *ImageBuilderRecipe) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("arn", e.arn)
	return properties
}

func (e *ImageBuilderRecipe) String() string {
	return e.arn
}
