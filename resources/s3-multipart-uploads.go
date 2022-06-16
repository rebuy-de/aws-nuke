package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type S3MultipartUpload struct {
	svc      *s3.S3
	bucket   string
	key      string
	uploadID string
}

func init() {
	register("S3MultipartUpload", ListS3MultipartUpload)
}

func ListS3MultipartUpload(sess *session.Session) ([]Resource, error) {
	svc := s3.New(sess)

	resources := make([]Resource, 0)

	buckets, err := DescribeS3Buckets(svc)
	if err != nil {
		return nil, err
	}

	for _, bucket := range buckets {
		params := &s3.ListMultipartUploadsInput{
			Bucket: bucket.Name,
		}

		for {
			resp, err := svc.ListMultipartUploads(params)
			if err != nil {
				return nil, err
			}

			for _, upload := range resp.Uploads {
				if upload.Key == nil || upload.UploadId == nil {
					continue
				}

				resources = append(resources, &S3MultipartUpload{
					svc:      svc,
					bucket:   aws.StringValue(bucket.Name),
					key:      *upload.Key,
					uploadID: *upload.UploadId,
				})
			}

			if *resp.IsTruncated {
				params.KeyMarker = resp.NextKeyMarker
				continue
			}

			break
		}
	}

	return resources, nil
}

func (e *S3MultipartUpload) Remove() error {
	params := &s3.AbortMultipartUploadInput{
		Bucket:   &e.bucket,
		Key:      &e.key,
		UploadId: &e.uploadID,
	}

	_, err := e.svc.AbortMultipartUpload(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *S3MultipartUpload) Properties() types.Properties {
	return types.NewProperties().
		Set("Bucket", e.bucket).
		Set("Key", e.key).
		Set("UploadID", e.uploadID)
}

func (e *S3MultipartUpload) String() string {
	return fmt.Sprintf("s3://%s/%s#%s", e.bucket, e.key, e.uploadID)
}
