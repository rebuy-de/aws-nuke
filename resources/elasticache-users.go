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

	params := &elasticache.DescribeUsersInput{MaxRecords: aws.Int64(100)}
	resp, err := svc.DescribeUsers(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, user := range resp.Users {
		resources = append(resources, &ElasticacheUser{
			svc:      svc,
			userId:   user.UserId,
			userName: user.UserName,
		})

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
