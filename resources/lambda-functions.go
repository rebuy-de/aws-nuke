package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaFunction struct {
	svc          *lambda.Lambda
	functionName *string
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
		resources = append(resources, &LambdaFunction{
			svc:          svc,
			functionName: function.FunctionName,
		})
	}

	return resources, nil
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
