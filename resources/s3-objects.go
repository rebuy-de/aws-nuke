package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type S3Object struct {
	svc          *s3.S3
	bucket       string
	creationDate time.Time
	key          string
	versionID    *string
	latest       bool
}

func init() {
	register("S3Object", ListS3Objects)
}

func ListS3Objects(sess *session.Session) ([]Resource, error) {
	svc := s3.New(sess)

	resources := make([]Resource, 0)

	buckets, err := DescribeS3Buckets(svc)
	if err != nil {
		return nil, err
	}

	for _, bucket := range buckets {
		params := &s3.ListObjectVersionsInput{
			Bucket: bucket.Name,
		}

		for {
			resp, err := svc.ListObjectVersions(params)
			if err != nil {
				return nil, err
			}

			for _, out := range resp.Versions {
				if out.Key == nil {
					continue
				}

				resources = append(resources, &S3Object{
					svc:          svc,
					bucket:       aws.StringValue(bucket.Name),
					creationDate: aws.TimeValue(bucket.CreationDate),
					key:          *out.Key,
					versionID:    out.VersionId,
					latest:       UnPtrBool(out.IsLatest, false),
				})
			}

			for _, out := range resp.DeleteMarkers {
				if out.Key == nil {
					continue
				}

				resources = append(resources, &S3Object{
					svc:          svc,
					bucket:       aws.StringValue(bucket.Name),
					creationDate: aws.TimeValue(bucket.CreationDate),
					key:          *out.Key,
					versionID:    out.VersionId,
					latest:       UnPtrBool(out.IsLatest, false),
				})
			}

			// make sure to list all with more than 1000 objects
			if *resp.IsTruncated {
				params.KeyMarker = resp.NextKeyMarker
				continue
			}

			break
		}
	}

	return resources, nil
}

func (e *S3Object) Remove() error {
	params := &s3.DeleteObjectInput{
		Bucket:    &e.bucket,
		Key:       &e.key,
		VersionId: e.versionID,
	}

	_, err := e.svc.DeleteObject(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *S3Object) Properties() types.Properties {
	return types.NewProperties().
		Set("Bucket", e.bucket).
		Set("Key", e.key).
		Set("VersionID", e.versionID).
		Set("IsLatest", e.latest).
		Set("CreationDate", e.creationDate)
}

func (e *S3Object) String() string {
	if e.versionID != nil && *e.versionID != "null" && !e.latest {
		return fmt.Sprintf("s3://%s/%s#%s", e.bucket, e.key, *e.versionID)
	}
	return fmt.Sprintf("s3://%s/%s", e.bucket, e.key)
}
