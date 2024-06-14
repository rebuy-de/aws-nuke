package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sirupsen/logrus"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CognitoUserPoolDomain struct {
	svc          *cognitoidentityprovider.CognitoIdentityProvider
	name         *string
	userPoolName *string
	userPoolId   *string
	userPoolArn  *string
}

func init() {
	register("CognitoUserPoolDomain", ListCognitoUserPoolDomains)
}

func ListCognitoUserPoolDomains(sess *session.Session) ([]Resource, error) {
	svc := cognitoidentityprovider.New(sess)

	// Lookup current account ID
	stsSvc := sts.New(sess)
	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
	   return nil, err
	}
	accountId := callerID.Account
	region := sess.Config.Region

	userPools, poolErr := ListCognitoUserPools(sess)
	if poolErr != nil {
		return nil, poolErr
	}

	resources := make([]Resource, 0)
	for _, userPoolResource := range userPools {
		userPool, ok := userPoolResource.(*CognitoUserPool)
		if !ok {
			logrus.Errorf("Unable to case CognitoUserPool")
			continue
		}

		describeParams := &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: userPool.id,
		}
		userPoolDetails, err := svc.DescribeUserPool(describeParams)
		if err != nil {
			return nil, err
		}
		if userPoolDetails.UserPool.Domain == nil {
			// No domain on this user pool so skip
			continue
		}

		resources = append(resources, &CognitoUserPoolDomain{
			svc:          svc,
			name:         userPoolDetails.UserPool.Domain,
			userPoolName: userPool.name,
			userPoolId:   userPool.id,
			userPoolArn:  aws.String("arn:aws:cognito-idp:" + *region + ":" + *accountId + ":userpool/" + *userPool.id),
		})
	}

	return resources, nil
}

func (f *CognitoUserPoolDomain) Remove() error {
	params := &cognitoidentityprovider.DeleteUserPoolDomainInput{
		Domain:     f.name,
		UserPoolId: f.userPoolId,
	}
	_, err := f.svc.DeleteUserPoolDomain(params)

	return err
}

func (f *CognitoUserPoolDomain) Properties() types.Properties {
	properties := types.NewProperties()
	params := &cognitoidentityprovider.ListTagsForResourceInput{
		ResourceArn: f.userPoolArn,
	}
	tags, _ := f.svc.ListTagsForResource(params)
	// Get the tags from CognitoUserPool instead because CognitoUserPoolDomain
	// doesnt support tags and could get it from main resource which is CognitoUserPool
	for tagKey, tagValue := range tags.Tags {
		properties.SetTag(&tagKey, tagValue)
	}
	properties.Set("name", f.name)
	properties.Set("userPoolArn", f.userPoolArn)
	properties.Set("userPoolName", f.userPoolName)
	properties.Set("userPoolId", f.userPoolId)
	return properties
}

func (f *CognitoUserPoolDomain) String() string {
	return *f.userPoolName + " -> " + *f.name
}
