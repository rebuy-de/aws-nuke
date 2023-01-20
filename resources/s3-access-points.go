package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("S3AccessPoint", ListS3AccessPoints)
}

type S3AccessPoint struct {
	svc         *s3control.S3Control
	accountId   *string
	accessPoint *s3control.AccessPoint
}

func ListS3AccessPoints(s *session.Session) ([]Resource, error) {
	// Lookup current account ID
	stsSvc := sts.New(s)
	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	accountId := callerID.Account

	resources := []Resource{}
	svc := s3control.New(s)
	for {
		params := &s3control.ListAccessPointsInput{
			AccountId: accountId,
		}

		resp, err := svc.ListAccessPoints(params)
		if err != nil {
			return nil, err
		}

		for _, accessPoint := range resp.AccessPointList {
			resources = append(resources, &S3AccessPoint{
				svc:         svc,
				accountId:   accountId,
				accessPoint: accessPoint,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (e *S3AccessPoint) Remove() error {
	_, err := e.svc.DeleteAccessPoint(&s3control.DeleteAccessPointInput{
		AccountId: e.accountId,
		Name:      aws.String(*e.accessPoint.Name),
	})
	return err
}

func (e *S3AccessPoint) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("AccessPointArn", e.accessPoint.AccessPointArn).
		Set("Alias", e.accessPoint.Alias).
		Set("Bucket", e.accessPoint.Bucket).
		Set("Name", e.accessPoint.Name).
		Set("NetworkOrigin", e.accessPoint.NetworkOrigin)

	return properties
}

func (e *S3AccessPoint) String() string {
	return *e.accessPoint.AccessPointArn
}
