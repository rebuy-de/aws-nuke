package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type LightsailKeyPair struct {
	svc         *lightsail.Lightsail
	keyPairName *string
}

func init() {
	register("LightsailKeyPair", ListLightsailKeyPairs)
}

func ListLightsailKeyPairs(sess *session.Session) ([]Resource, error) {
	svc := lightsail.New(sess)
	resources := []Resource{}

	params := &lightsail.GetKeyPairsInput{}

	for {
		output, err := svc.GetKeyPairs(params)
		if err != nil {
			return nil, err
		}

		for _, keyPair := range output.KeyPairs {
			resources = append(resources, &LightsailKeyPair{
				svc:         svc,
				keyPairName: keyPair.Name,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

func (f *LightsailKeyPair) Remove() error {

	_, err := f.svc.DeleteKeyPair(&lightsail.DeleteKeyPairInput{
		KeyPairName: f.keyPairName,
	})

	return err
}

func (f *LightsailKeyPair) String() string {
	return *f.keyPairName
}
