package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type LambdaFunction struct {
	svc          *lambda.Lambda
	functionName *string
	tags         *map[string]*string
}

func init() {
	register("LambdaFunction", ListLambdaFunctions)
}

func ListLambdaFunctions(sess *session.Session) ([]Resource, error) {
	svc := lambda.New(sess)

	params := &lambda.ListFunctionsInput{}
	resp, err := svc.ListFunctions(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, function := range resp.Functions {
		tags, err := retrieveLambdaTags(svc, *function.FunctionArn)
		if err != nil {
			continue
		}

		resources = append(resources, &LambdaFunction{
			svc:          svc,
			functionName: function.FunctionName,
			tags: &tags,
		})
	}

	return resources, nil
}

func retrieveLambdaTags(svc *lambda.Lambda, arn string) (map[string]*string, error) {
	input := &lambda.ListTagsInput {
		Resource: aws.String(arn),
	}

	result, err := svc.ListTags(input)

	if err != nil {
		return make(map[string]*string, 0), err
	}

	return result.Tags, nil
}

func (f *LambdaFunction) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.functionName)

	for key, val := range *f.tags {
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
