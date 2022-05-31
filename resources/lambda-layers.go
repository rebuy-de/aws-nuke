package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type lambdaLayer struct {
	svc       *lambda.Lambda
	layerName *string
	version   int64
}

func init() {
	register("LambdaLayer", ListLambdaLayers)
}

func ListLambdaLayers(sess *session.Session) ([]Resource, error) {
	svc := lambda.New(sess)

	layers := make([]*lambda.LayersListItem, 0)

	params := &lambda.ListLayersInput{}

	err := svc.ListLayersPages(params, func(page *lambda.ListLayersOutput, lastPage bool) bool {
		for _, out := range page.Layers {
			layers = append(layers, out)
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, layer := range layers {
		versionsParams := &lambda.ListLayerVersionsInput{
			LayerName: layer.LayerName,
		}
		err := svc.ListLayerVersionsPages(versionsParams, func(page *lambda.ListLayerVersionsOutput, lastPage bool) bool {
			for _, out := range page.LayerVersions {
				resources = append(resources, &lambdaLayer{
					svc:       svc,
					layerName: layer.LayerName,
					version:   *out.Version,
				})
			}
			return true
		})
		if err != nil {
			return nil, err
		}

	}

	return resources, nil
}

func (l *lambdaLayer) Remove() error {

	_, err := l.svc.DeleteLayerVersion(&lambda.DeleteLayerVersionInput{
		LayerName:     l.layerName,
		VersionNumber: &l.version,
	})

	return err
}

func (l *lambdaLayer) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", l.layerName)
	properties.Set("Version", l.version)

	return properties
}

func (l *lambdaLayer) String() string {
	return *l.layerName
}
