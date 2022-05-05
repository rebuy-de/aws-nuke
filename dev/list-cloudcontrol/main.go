package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/rebuy-de/aws-nuke/v2/resources"
	"github.com/rebuy-de/rebuy-go-sdk/v3/pkg/cmdutil"
	"github.com/sirupsen/logrus"
)

type CFTypeSchema struct {
	Handlers map[string]interface{} `json:"handlers"`
}

func main() {
	ctx := cmdutil.SignalRootContext()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	cf := cloudformation.New(sess)

	mapping := resources.GetCloudControlMapping()

	in := &cloudformation.ListTypesInput{
		Type:             aws.String(cloudformation.RegistryTypeResource),
		Visibility:       aws.String(cloudformation.VisibilityPublic),
		ProvisioningType: aws.String(cloudformation.ProvisioningTypeFullyMutable),
	}

	err = cf.ListTypesPagesWithContext(ctx, in, func(out *cloudformation.ListTypesOutput, _ bool) bool {
		if out == nil {
			return true
		}

		for _, summary := range out.TypeSummaries {
			if summary == nil {
				continue
			}

			typeName := aws.StringValue(summary.TypeName)
			color.New(color.Bold).Printf("%-55s", typeName)
			if !strings.HasPrefix(typeName, "AWS::") {
				color.HiBlack("does not have a valid prefix")
				continue
			}

			describe, err := cf.DescribeType(&cloudformation.DescribeTypeInput{
				Type:     aws.String(cloudformation.RegistryTypeResource),
				TypeName: aws.String(typeName),
			})
			if err != nil {
				color.New(color.FgRed).Println(err)
				continue
			}

			var schema CFTypeSchema
			err = json.Unmarshal([]byte(aws.StringValue(describe.Schema)), &schema)
			if err != nil {
				color.New(color.FgRed).Println(err)
				continue
			}

			_, canList := schema.Handlers["list"]
			if !canList {
				color.New(color.FgHiBlack).Println("does not support list")
				continue
			}

			resourceName, exists := mapping[typeName]
			if exists && resourceName == typeName {
				fmt.Print("is only covered by ")
				color.New(color.FgGreen, color.Bold).Println(resourceName)
				continue
			} else if exists {
				fmt.Print("is also covered by ")
				color.New(color.FgBlue, color.Bold).Println(resourceName)
				continue
			}

			color.New(color.FgYellow).Println("is not configured")
		}

		return true
	})
	if err != nil {
		logrus.Fatal(err)
	}
}
