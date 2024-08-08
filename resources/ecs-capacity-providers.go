package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ECSCapacityProvider struct {
	svc *ecs.ECS
	ARN *string
}

func init() {
	register("ECSCapacityProvider", DescribeECSCapacityProviders)
}

func DescribeECSCapacityProviders(sess *session.Session) ([]Resource, error) {
	svc := ecs.New(sess)
	resources := []Resource{}

	params := &ecs.DescribeCapacityProvidersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeCapacityProviders(params)
		if err != nil {
			return nil, err
		}

		for _, capacityProviders := range output.CapacityProviders {
			if *capacityProviders.Name == "FARGATE" || *capacityProviders.Name == "FARGATE_SPOT" {
				// The FARGATE and FARGATE_SPOT capacity providers cannot be deleted
				continue
			}

			resources = append(resources, &ECSCapacityProvider{
				svc: svc,
				ARN: capacityProviders.CapacityProviderArn,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ECSCapacityProvider) Remove() error {

	_, err := f.svc.DeleteCapacityProvider(&ecs.DeleteCapacityProviderInput{
		CapacityProvider: f.ARN,
	})

	return err
}

func (f *ECSCapacityProvider) String() string {
	return *f.ARN
}
