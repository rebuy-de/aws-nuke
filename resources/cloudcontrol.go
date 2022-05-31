package resources

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudcontrolapi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

func init() {
	// It is required to manually define Cloud Control API targets, because
	// existing configs that already filter old-style resources could break,
	// because the resource is also available via Cloud Control.
	//
	// To get an overview of available cloud control resource types run this
	// command in the repo root:
	//     go run ./dev/list-cloudcontrol
	registerCloudControl("AWS::AppFlow::ConnectorProfile")
	registerCloudControl("AWS::AppFlow::Flow")
	registerCloudControl("AWS::AppRunner::Service")
	registerCloudControl("AWS::ApplicationInsights::Application")
	registerCloudControl("AWS::Backup::Framework")
	registerCloudControl("AWS::MWAA::Environment")
	registerCloudControl("AWS::Synthetics::Canary")
	registerCloudControl("AWS::Timestream::Database")
	registerCloudControl("AWS::Timestream::ScheduledQuery")
	registerCloudControl("AWS::Timestream::Table")
	registerCloudControl("AWS::Transfer::Workflow")
	registerCloudControl("AWS::NetworkFirewall::Firewall")
	registerCloudControl("AWS::NetworkFirewall::FirewallPolicy")
	registerCloudControl("AWS::NetworkFirewall::RuleGroup")
}

func NewListCloudControlResource(typeName string) func(*session.Session) ([]Resource, error) {
	return func(sess *session.Session) ([]Resource, error) {
		svc := cloudcontrolapi.New(sess)

		params := &cloudcontrolapi.ListResourcesInput{
			TypeName: aws.String(typeName),
		}
		resources := make([]Resource, 0)
		err := svc.ListResourcesPages(params, func(page *cloudcontrolapi.ListResourcesOutput, lastPage bool) bool {
			for _, desc := range page.ResourceDescriptions {
				identifier := aws.StringValue(desc.Identifier)

				properties, err := cloudControlParseProperties(aws.StringValue(desc.Properties))
				if err != nil {
					logrus.
						WithError(errors.WithStack(err)).
						WithField("type-name", typeName).
						WithField("identifier", identifier).
						Error("failed to parse cloud control properties")
					continue
				}
				properties = properties.Set("Identifier", identifier)
				resources = append(resources, &CloudControlResource{
					svc:         svc,
					clientToken: uuid.New().String(),
					typeName:    typeName,
					identifier:  identifier,
					properties:  properties,
				})
			}

			return true
		})

		if err != nil {
			return nil, err
		}

		return resources, nil
	}
}

func cloudControlParseProperties(payload string) (types.Properties, error) {
	// Warning: The implementation of this function is not very straighforward,
	// because the aws-nuke filter functions expect a very rigid structure and
	// the properties from the Cloud Control API are very dynamic.

	propMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(payload), &propMap)
	if err != nil {
		return nil, err
	}
	properties := types.NewProperties()
	for name, value := range propMap {
		switch v := value.(type) {
		case string:
			properties = properties.Set(name, v)
		case []interface{}:
			for _, value2 := range v {
				switch v2 := value2.(type) {
				case string:
					properties.Set(
						fmt.Sprintf("%s.[%q]", name, v2),
						true,
					)
				case map[string]interface{}:
					if len(v2) == 2 && v2["Key"] != nil && v2["Value"] != nil {
						properties.Set(
							fmt.Sprintf("%s.[%q]", name, v2["Key"]),
							v2["Value"],
						)
					} else {
						logrus.
							WithField("value", fmt.Sprintf("%q", v)).
							Debugf("nested cloud control property type []%T is not supported", value)
					}
				default:
					logrus.
						WithField("value", fmt.Sprintf("%q", v)).
						Debugf("nested cloud control property type []%T is not supported", value)
				}
			}

		default:
			// We cannot rely on the default handling of
			// properties.Set, because it would fall back to
			// fmt.Sprintf. Since the cloud control properties are
			// nested it would create properties that are not
			// suitable for filtering. Therefore we have to
			// implemented more sophisticated parsing.
			logrus.
				WithField("value", fmt.Sprintf("%q", v)).
				Debugf("cloud control property type %T is not supported", v)
		}
	}

	return properties, nil
}

type CloudControlResource struct {
	svc         *cloudcontrolapi.CloudControlApi
	clientToken string
	typeName    string
	identifier  string
	properties  types.Properties
}

func (r *CloudControlResource) String() string {
	return r.identifier
}

func (i *CloudControlResource) Remove() error {
	_, err := i.svc.DeleteResource(&cloudcontrolapi.DeleteResourceInput{
		ClientToken: &i.clientToken,
		Identifier:  &i.identifier,
		TypeName:    &i.typeName,
	})
	return err
}

func (r *CloudControlResource) Properties() types.Properties {
	return r.properties
}
