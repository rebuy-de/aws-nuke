package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CognitoUserPool struct {
	svc  *cognitoidentityprovider.CognitoIdentityProvider
	name *string
	id   *string
	arn  *string
}

func init() {
	register("CognitoUserPool", ListCognitoUserPools)
}

func ListCognitoUserPools(sess *session.Session) ([]Resource, error) {
	svc := cognitoidentityprovider.New(sess)

	// Lookup current account ID
	stsSvc := sts.New(sess)
	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		   return nil, err
	}
	accountId := callerID.Account
	region := sess.Config.Region

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
				arn:  aws.String("arn:aws:cognito-idp:" + *region + ":" + *accountId + ":userpool/" + *pool.Id),
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

	_, err := f.svc.DeleteUserPool(&cognitoidentityprovider.DeleteUserPoolInput{
		UserPoolId: f.id,
	})

	return err
}

func (f *CognitoUserPool) Properties() types.Properties {
	properties := types.NewProperties()
	params := &cognitoidentityprovider.ListTagsForResourceInput{
		ResourceArn: f.arn,
	}
	tags, _ := f.svc.ListTagsForResource(params)
	for tagKey, tagValue := range tags.Tags {
		properties.SetTag(&tagKey, tagValue)
	}
	properties.Set("name", f.name)
	properties.Set("id", f.id)
	return properties
}

func (f *CognitoUserPool) String() string {
	return *f.name
}
