package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
)

type CognitoIdentityPool struct {
	svc  *cognitoidentity.CognitoIdentity
	name *string
	id   *string
}

func init() {
	register("CognitoIdentityPool", ListCognitoIdentityPools)
}

func ListCognitoIdentityPools(sess *session.Session) ([]Resource, error) {
	svc := cognitoidentity.New(sess)
	resources := []Resource{}

	params := &cognitoidentity.ListIdentityPoolsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListIdentityPools(params)
		if err != nil {
			return nil, err
		}

		for _, pool := range output.IdentityPools {
			resources = append(resources, &CognitoIdentityPool{
				svc:  svc,
				name: pool.IdentityPoolName,
				id:   pool.IdentityPoolId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CognitoIdentityPool) Remove() error {

	_, err := f.svc.DeleteIdentityPool(&cognitoidentity.DeleteIdentityPoolInput{
		IdentityPoolId: f.id,
	})

	return err
}

func (f *CognitoIdentityPool) String() string {
	return *f.name
}
