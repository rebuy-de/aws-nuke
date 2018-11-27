package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoUserPool struct {
	svc    *cognitoidentityprovider.CognitoIdentityProvider
	name   *string
	id     *string
	domain *string
}

func init() {
	register("CognitoUserPool", ListCognitoUserPools)
}

func ListCognitoUserPools(sess *session.Session) ([]Resource, error) {
	svc := cognitoidentityprovider.New(sess)
	resources := []Resource{}

	params := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListUserPools(params)
		if err != nil {
			return nil, err
		}

		for _, pool := range output.UserPools {
			describeParams := &cognitoidentityprovider.DescribeUserPoolInput{
				UserPoolId: pool.Id,
			}
			userpool, err := svc.DescribeUserPool(describeParams)
			if err != nil {
				return nil, err
			}
			domain := ""
			if userpool.UserPool.Domain != nil {
				domain = *userpool.UserPool.Domain
			}

			resources = append(resources, &CognitoUserPool{
				svc:    svc,
				name:   pool.Name,
				id:     pool.Id,
				domain: &domain,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CognitoUserPool) Remove() error {
	if *f.domain != "" {
		domainParams := &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     f.domain,
			UserPoolId: f.id,
		}
		_, err := f.svc.DeleteUserPoolDomain(domainParams)
		if err != nil {
			return err
		}
	}

	_, err := f.svc.DeleteUserPool(&cognitoidentityprovider.DeleteUserPoolInput{
		UserPoolId: f.id,
	})

	return err
}

func (f *CognitoUserPool) String() string {
	return *f.name
}
