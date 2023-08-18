package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
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

	var resources []Resource
	var marker *string // Marker for pagination

	for {
		params := &elasticache.DescribeUsersInput{
			MaxRecords: aws.Int64(100),
			Marker:     marker,
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

		// If there are more results, the response will have a Marker.
		// Set the marker for the next iteration.
		if resp.Marker != nil {
			marker = resp.Marker
		} else {
			// No more results, break the loop.
			break
		}
	}

	return resources, nil
}

func (i *ElasticacheUser) Filter() error {
	if strings.HasPrefix(*i.userName, "default") {
		return fmt.Errorf("Cannot delete default user")
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

func (i *ElasticacheUser) String() string {
	return *i.userId
}
