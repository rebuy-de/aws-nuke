package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transfer"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TransferServerUser struct {
	svc      *transfer.Transfer
	username *string
	serverID *string
	tags     []*transfer.Tag
}

func init() {
	register("TransferServerUser", ListTransferServerUsers)
}

func ListTransferServerUsers(sess *session.Session) ([]Resource, error) {
	svc := transfer.New(sess)
	resources := []Resource{}

	params := &transfer.ListServersInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListServers(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Servers {
			userParams := &transfer.ListUsersInput{
				MaxResults: aws.Int64(100),
				ServerId:   item.ServerId,
			}

			for {
				userOutput, err := svc.ListUsers(userParams)
				if err != nil {
					return nil, err
				}

				for _, user := range userOutput.Users {
					descOutput, err := svc.DescribeUser(&transfer.DescribeUserInput{
						ServerId: item.ServerId,
						UserName: user.UserName,
					})
					if err != nil {
						return nil, err
					}

					resources = append(resources, &TransferServerUser{
						svc:      svc,
						username: user.UserName,
						serverID: item.ServerId,
						tags:     descOutput.User.Tags,
					})

				}

				if userOutput.NextToken == nil {
					break
				}

				userParams.NextToken = userOutput.NextToken
			}
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (ts *TransferServerUser) Remove() error {

	_, err := ts.svc.DeleteUser(&transfer.DeleteUserInput{
		ServerId: ts.serverID,
		UserName: ts.username,
	})

	return err
}

func (ts *TransferServerUser) String() string {
	return fmt.Sprintf("%s -> %s", *ts.serverID, *ts.username)
}

func (ts *TransferServerUser) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range ts.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.
		Set("Username", ts.username).
		Set("ServerID", ts.serverID)
	return properties
}
