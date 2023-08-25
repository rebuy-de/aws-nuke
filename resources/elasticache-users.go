package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ElasticacheUser struct {
	svc      *elasticache.ElastiCache
	userId   *string
	userName *string
}

func init() {
	register("ElasticacheUser", ListElasticacheUsers)
}

func ListElasticacheUsers(sess *session.Session) ([]Resource, error) {
	svc := elasticache.New(sess)
	resources := []Resource{}
	var nextToken *string

	for {
		params := &elasticache.DescribeUsersInput{
			MaxRecords: aws.Int64(100),
			Marker:     nextToken,
		}
		resp, err := svc.DescribeUsers(params)
		if err != nil {
			return nil, err
		}

		for _, user := range resp.Users {
			resources = append(resources, &ElasticacheUser{
				svc:      svc,
				userId:   user.UserId,
				userName: user.UserName,
			})
		}

		// Check if there are more results
		if resp.Marker == nil {
			break // No more results, exit the loop
		}

		// Set the nextToken for the next iteration
		nextToken = resp.Marker
	}

	return resources, nil
}

func (i *ElasticacheUser) Filter() error {
	if strings.HasPrefix(*i.userName, "default") {
		return fmt.Errorf("cannot delete default user")
	}
	return nil
}

func (i *ElasticacheUser) Remove() error {
	params := &elasticache.DeleteUserInput{
		UserId: i.userId,
	}

	_, err := i.svc.DeleteUser(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *ElasticacheUser) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", i.userId)
	properties.Set("UserName", i.userName)
	return properties
}

func (i *ElasticacheUser) String() string {
	return *i.userId
}
