package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type CognitoIdentityProvider struct {
	svc          *cognitoidentityprovider.CognitoIdentityProvider
	name         *string
	providerType *string
	userPoolName *string
	userPoolId   *string
}

func init() {
	register("CognitoIdentityProvider", ListCognitoIdentityProviders)
}

func ListCognitoIdentityProviders(sess *session.Session) ([]Resource, error) {
	svc := cognitoidentityprovider.New(sess)

	userPools, poolErr := ListCognitoUserPools(sess)
	if poolErr != nil {
		return nil, poolErr
	}

	resources := []Resource{}

	for _, userPoolResource := range userPools {
		userPool, ok := userPoolResource.(*CognitoUserPool)
		if !ok {
			logrus.Errorf("Unable to case CognitoUserPool")
			continue
		}

		listParams := &cognitoidentityprovider.ListIdentityProvidersInput{
			UserPoolId: userPool.id,
			MaxResults: aws.Int64(50),
		}

		for {
			output, err := svc.ListIdentityProviders(listParams)
			if err != nil {
				return nil, err
			}

			for _, provider := range output.Providers {
				resources = append(resources, &CognitoIdentityProvider{
					svc:          svc,
					name:         provider.ProviderName,
					providerType: provider.ProviderType,
					userPoolName: userPool.name,
					userPoolId:   userPool.id,
				})
			}

			if output.NextToken == nil {
				break
			}

			listParams.NextToken = output.NextToken
		}
	}

	return resources, nil
}

func (p *CognitoIdentityProvider) Remove() error {

	_, err := p.svc.DeleteIdentityProvider(&cognitoidentityprovider.DeleteIdentityProviderInput{
		UserPoolId:   p.userPoolId,
		ProviderName: p.name,
	})

	return err
}

func (p *CognitoIdentityProvider) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Type", p.providerType)
	properties.Set("UserPoolName", p.userPoolName)
	properties.Set("Name", p.name)
	return properties
}

func (p *CognitoIdentityProvider) String() string {
	return *p.userPoolName + " -> " + *p.name
}
