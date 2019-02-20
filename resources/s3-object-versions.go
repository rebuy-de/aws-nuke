package resources

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DeleteListVersionListIterator is an iterator for S3 Object Versions.
type DeleteVersionListIterator struct {
	Bucket    *string
	Paginator request.Pagination
	objects   []*s3.ObjectVersion
}

// NewDeleteListVersionListIterator will return a new BatchDeleteIterator for S3 Object Versions.
func NewDeleteVersionListIterator(svc s3iface.S3API, input *s3.ListObjectVersionsInput, opts ...func(*DeleteVersionListIterator)) s3manager.BatchDeleteIterator {
	iter := &DeleteVersionListIterator{
		Bucket: input.Bucket,
		Paginator: request.Pagination{
			NewRequest: func() (*request.Request, error) {
				var inCpy *s3.ListObjectVersionsInput
				if input != nil {
					tmp := *input
					inCpy = &tmp
				}
				req, _ := svc.ListObjectVersionsRequest(inCpy)
				return req, nil
			},
		},
	}

	for _, opt := range opts {
		opt(iter)
	}
	return iter
}

// Next will use the S3API client to iterate through a list of objects.
func (iter *DeleteVersionListIterator) Next() bool {
	if len(iter.objects) > 0 {
		iter.objects = iter.objects[1:]
	}

	if len(iter.objects) == 0 && iter.Paginator.Next() {
		output := iter.Paginator.Page().(*s3.ListObjectVersionsOutput)
		iter.objects = output.Versions

		for _, entry := range output.DeleteMarkers {
			iter.objects = append(iter.objects, &s3.ObjectVersion{
				Key:       entry.Key,
				VersionId: entry.VersionId,
			})
		}
	}

	return len(iter.objects) > 0
}

// Err will return the last known error from Next.
func (iter *DeleteVersionListIterator) Err() error {
	return iter.Paginator.Err()
}

// DeleteObject will return the current object to be deleted.
func (iter *DeleteVersionListIterator) DeleteObject() s3manager.BatchDeleteObject {
	return s3manager.BatchDeleteObject{
		Object: &s3.DeleteObjectInput{
			Bucket:    iter.Bucket,
			Key:       iter.objects[0].Key,
			VersionId: iter.objects[0].VersionId,
		},
	}
}
