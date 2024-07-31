package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"strings"
	"time"
)

type FSxFileSystem struct {
	svc        *fsx.FSx
	filesystem *fsx.FileSystem
}

func init() {
	register("FSxFileSystem", ListFSxFileSystems)
}

func ListFSxFileSystems(sess *session.Session) ([]Resource, error) {
	svc := fsx.New(sess)
	resources := []Resource{}

	params := &fsx.DescribeFileSystemsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.DescribeFileSystems(params)
		if err != nil {
			return nil, err
		}

		for _, filesystem := range resp.FileSystems {
			resources = append(resources, &FSxFileSystem{
				svc:        svc,
				filesystem: filesystem,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func isRootVolume(volume *fsx.Volume) bool {
	return strings.HasSuffix(*volume.Name, "_root")
}

func deleteFileSystemWithRetry(svc *fsx.FSx, filesystemId *string) error {
	const maxRetries = 5
	const retryDelay = 10 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		describeSVMsInput := &fsx.DescribeStorageVirtualMachinesInput{
			Filters: []*fsx.StorageVirtualMachineFilter{
				{
					Name:   aws.String("file-system-id"),
					Values: []*string{filesystemId},
				},
			},
		}
		svmResp, err := svc.DescribeStorageVirtualMachines(describeSVMsInput)
		if err != nil {
			return err
		}

		if len(svmResp.StorageVirtualMachines) == 0 {
			_, err := svc.DeleteFileSystem(&fsx.DeleteFileSystemInput{
				FileSystemId: filesystemId,
			})
			if err != nil {
				return err
			}
			return nil
		} else {
			time.Sleep(retryDelay)
		}
	}

	return awserr.New("MaxRetriesExceeded", "unable to delete filesystem after max retries", nil)
}

func (f *FSxFileSystem) Remove() error {
	describeSVMsInput := &fsx.DescribeStorageVirtualMachinesInput{
		StorageVirtualMachineIds: []*string{ /* Populate with SVM IDs as needed */ },
	}

	svmResp, err := f.svc.DescribeStorageVirtualMachines(describeSVMsInput)
	if err != nil {
		return err
	}

	for _, svm := range svmResp.StorageVirtualMachines {
		listVolumesInput := &fsx.DescribeVolumesInput{
			Filters: []*fsx.VolumeFilter{
				{
					Name:   aws.String("storage-virtual-machine-id"),
					Values: []*string{svm.StorageVirtualMachineId},
				},
			},
		}

		volumeResp, err := f.svc.DescribeVolumes(listVolumesInput)
		if err != nil {
			return err
		}

		for _, volume := range volumeResp.Volumes {
			if !isRootVolume(volume) {
				_, err := f.svc.DeleteVolume(&fsx.DeleteVolumeInput{
					VolumeId: volume.VolumeId,
					OntapConfiguration: &fsx.DeleteVolumeOntapConfiguration{
						SkipFinalBackup: aws.Bool(true),
					},
				})
				if err != nil {
					return err
				}
			}
		}

		_, err = f.svc.DeleteStorageVirtualMachine(&fsx.DeleteStorageVirtualMachineInput{
			StorageVirtualMachineId: svm.StorageVirtualMachineId,
		})

		if err != nil {
			return err
		}
	}

	fsDeletionErr := deleteFileSystemWithRetry(f.svc, f.filesystem.FileSystemId)
	if fsDeletionErr != nil {
		return fsDeletionErr
	}

	return nil
}

func (f *FSxFileSystem) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range f.filesystem.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("Type", f.filesystem.FileSystemType)
	return properties
}

func (f *FSxFileSystem) String() string {
	return *f.filesystem.FileSystemId
}
