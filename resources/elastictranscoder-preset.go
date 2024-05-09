package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ElasticTranscoderPreset struct {
	svc      *elastictranscoder.ElasticTranscoder
	presetID *string
}

func init() {
	register("ElasticTranscoderPreset", ListElasticTranscoderPresets)
}

func ListElasticTranscoderPresets(sess *session.Session) ([]Resource, error) {
	svc := elastictranscoder.New(sess)
	resources := []Resource{}

	params := &elastictranscoder.ListPresetsInput{}

	for {
		resp, err := svc.ListPresets(params)
		if err != nil {
			return nil, err
		}

		for _, preset := range resp.Presets {
			resources = append(resources, &ElasticTranscoderPreset{
				svc:      svc,
				presetID: preset.Id,
			})
		}

		if resp.NextPageToken == nil {
			break
		}

		params.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ElasticTranscoderPreset) Filter() error {
	if strings.HasPrefix(*f.presetID, "1351620000001") {
		return fmt.Errorf("cannot delete elastic transcoder system presets")
	}
	return nil
}

func (f *ElasticTranscoderPreset) Remove() error {

	_, err := f.svc.DeletePreset(&elastictranscoder.DeletePresetInput{
		Id: f.presetID,
	})

	return err
}

func (f *ElasticTranscoderPreset) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("PresetID", f.presetID)
	return properties
}

func (f *ElasticTranscoderPreset) String() string {
	return *f.presetID
}
