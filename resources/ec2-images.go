package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Image struct {
	svc *ec2.EC2
	id  string
}

func init() {
	register("EC2Image", ListEC2Images)
}

func ListEC2Images(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeImagesInput{
		Owners: []*string{
			aws.String("self"),
		},
	}
	resp, err := svc.DescribeImages(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Images {
		resources = append(resources, &EC2Image{
			svc: svc,
			id:  *out.ImageId,
		})
	}

	return resources, nil
}

func (e *EC2Image) Remove() error {
	_, err := e.svc.DeregisterImage(&ec2.DeregisterImageInput{
		ImageId: &e.id,
	})
	return err
}

func (e *EC2Image) String() string {
	return e.id
}
