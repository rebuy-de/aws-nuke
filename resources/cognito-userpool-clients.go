package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type CognitoUserPoolClient struct {
	svc          *cognitoidentityprovider.CognitoIdentityProvider
	name         *string
	id           *string
	userPoolName *string
	userPoolId   *string
}

func init() {
	register("CognitoUserPoolClient", ListCognitoUserPoolClients)
}

func ListCognitoUserPoolClients(sess *session.Session) ([]Resource, error) {
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

		listParams := &cognitoidentityprovider.ListUserPoolClientsInput{
			UserPoolId: userPool.id,
			MaxResults: aws.Int64(50),
		}

		for {
			output, err := svc.ListUserPoolClients(listParams)
			if err != nil {
				return nil, err
			}

			for _, client := range output.UserPoolClients {
				resources = append(resources, &CognitoUserPoolClient{
					svc:          svc,
					id:           client.ClientId,
					name:         client.ClientName,
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

func (p *CognitoUserPoolClient) Remove() error {

	_, err := p.svc.DeleteUserPoolClient(&cognitoidentityprovider.DeleteUserPoolClientInput{
		ClientId:   p.id,
		UserPoolId: p.userPoolId,
	})

	return err
}

func (p *CognitoUserPoolClient) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", p.id)
	properties.Set("Name", p.name)
	properties.Set("UserPoolName", p.userPoolName)
	return properties
}

func (p *CognitoUserPoolClient) String() string {
	return *p.userPoolName + " -> " + *p.name
}
