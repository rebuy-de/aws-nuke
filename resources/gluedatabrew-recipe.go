package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gluedatabrew"
)

type GlueDataBrewRecipe struct {
	svc           *gluedatabrew.GlueDataBrew
	name          *string
	recipeversion *string
}

func init() {
	register("GlueDataBrewRecipe", ListGlueDataBrewRecipe)
}

func ListGlueDataBrewRecipe(sess *session.Session) ([]Resource, error) {
	svc := gluedatabrew.New(sess)
	resources := []Resource{}

	params := &gluedatabrew.ListRecipesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListRecipes(params)
		if err != nil {
			return nil, err
		}

		for _, recipe := range output.Recipes {
			resources = append(resources, &GlueDataBrewRecipe{
				svc:           svc,
				name:          recipe.Name,
				recipeversion: recipe.RecipeVersion,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDataBrewRecipe) Remove() error {
	_, err := f.svc.DeleteRecipeVersion(&gluedatabrew.DeleteRecipeVersionInput{
		Name:          f.name,
		RecipeVersion: f.recipeversion,
	})

	return err
}

func (f *GlueDataBrewRecipe) String() string {
	return *f.name
}
