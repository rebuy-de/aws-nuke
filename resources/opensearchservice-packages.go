package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type OSPackage struct {
	svc         *opensearchservice.OpenSearchService
	packageID   *string
	packageName *string
	createdTime *time.Time
}

func init() {
	register("OSPackage", ListOSPackages)
}

func ListOSPackages(sess *session.Session) ([]Resource, error) {
	svc := opensearchservice.New(sess)

	listResp, err := svc.DescribePackages(&opensearchservice.DescribePackagesInput{})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, pkg := range listResp.PackageDetailsList {
		resources = append(resources, &OSPackage{
			svc:         svc,
			packageID:   pkg.PackageID,
			packageName: pkg.PackageName,
			createdTime: pkg.CreatedAt,
		})
	}

	return resources, nil
}

func (o *OSPackage) Remove() error {
	_, err := o.svc.DeletePackage(&opensearchservice.DeletePackageInput{
		PackageID: o.packageID,
	})

	return err
}

func (o *OSPackage) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("PackageID", o.packageID)
	properties.Set("PackageName", o.packageName)
	properties.Set("CreatedTime", o.createdTime.Format(time.RFC3339))
	return properties
}

func (o *OSPackage) String() string {
	return *o.packageID
}
