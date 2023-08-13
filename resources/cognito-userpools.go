package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
)

type CognitoUserPool struct {
	svc          *cognitoidentityprovider.CognitoIdentityProvider
	name         *string
	id           *string
	featureFlags config.FeatureFlags
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
			resources = append(resources, &CognitoUserPool{
				svc:  svc,
				name: pool.Name,
				id:   pool.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (l *CognitoUserPool) FeatureFlags(ff config.FeatureFlags) {
	l.featureFlags = ff
}

func (f *CognitoUserPool) Remove() error {
	_, err := f.svc.DeleteUserPool(&cognitoidentityprovider.DeleteUserPoolInput{
		UserPoolId: f.id,
	})
	if err != nil {
		if f.featureFlags.DisableDeletionProtection.CognitoUserPool {
			err = f.DisableProtection()
			if err != nil {
				return err
			}
			_, err = f.svc.DeleteUserPool(&cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: f.id,
			})
			if err != nil {
				return err
			}
			return nil
		}
	}
	return err
}

func (e *CognitoUserPool) DisableProtection() error {
	userPoolOutput, err := e.svc.DescribeUserPool(&cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: e.id,
	})
	if err != nil {
		return err
	}
	userPool := userPoolOutput.UserPool
	params := &cognitoidentityprovider.UpdateUserPoolInput{
		DeletionProtection:     &cognitoidentityprovider.DeletionProtectionType_Values()[1],
		UserPoolId:             e.id,
		AutoVerifiedAttributes: userPool.AutoVerifiedAttributes,
	}
	_, err = e.svc.UpdateUserPool(params)
	return err
}

func (f *CognitoUserPool) String() string {
	return *f.name
}
