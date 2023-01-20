package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type LambdaFunction struct {
	svc          *lambda.Lambda
	functionName *string
	tags         map[string]*string
}

func init() {
	register("LambdaFunction", ListLambdaFunctions)
}

func ListLambdaFunctions(sess *session.Session) ([]Resource, error) {
	svc := lambda.New(sess)

	functions := make([]*lambda.FunctionConfiguration, 0)

	params := &lambda.ListFunctionsInput{}

	err := svc.ListFunctionsPages(params, func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
		for _, out := range page.Functions {
			functions = append(functions, out)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, function := range functions {
		tags, err := svc.ListTags(&lambda.ListTagsInput{
			Resource: function.FunctionArn,
		})

		if err != nil {
			continue
		}

		resources = append(resources, &LambdaFunction{
			svc:          svc,
			functionName: function.FunctionName,
			tags:         tags.Tags,
		})
	}

	return resources, nil
}

func (f *LambdaFunction) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.functionName)

	for key, val := range f.tags {
		properties.SetTag(&key, val)
	}

	return properties
}

func (f *LambdaFunction) Remove() error {

	_, err := f.svc.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: f.functionName,
	})

	return err
}

func (f *LambdaFunction) String() string {
	return *f.functionName
}
