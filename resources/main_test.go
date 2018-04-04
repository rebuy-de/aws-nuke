package resources

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var debugFlag = false

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	// Setup
	os.Setenv("TF_IN_AUTOMATION", "1")
	if err := RunCmd("terraform", "init", "-input=false", "-no-color"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Test
	code := m.Run()

	// Teardown
	fmt.Println("Teardown: destroying all terraform resources of the test suite")
	RunCmd("terraform", "destroy", "-input=false", "-no-color", "-auto-approve")

	os.Exit(code)
}

func ResourceTypeTest(tfType, tfName, nukeType string, t *testing.T) error {
	tfResourcePath := tfType + "." + tfName

	// Validate target resources do not exist yet.
	planfile, preCreateDiff, err := CreatePlan([]string{tfResourcePath})
	if len(preCreateDiff) == 0 {
		return fmt.Errorf(`no resources to create. The AWS resource may already exist or the terraform resource cannot be found: "%s"`, tfResourcePath)
	}
	fmt.Println("Plan to create the following resources:", preCreateDiff)

	// Apply the plan. Terraform also creates the resources' dependencies.
	if err := TerraformApplyPlan(planfile); err != nil {
		return err
	}

	// Some resources require other resources to be created. The terraform name
	// of these required resources start with "dep_".
	var requiredResources []string
	for _, r := range preCreateDiff {
		name := strings.Split(r, ".")[1]
		if strings.HasPrefix(name, "dep_") {
			requiredResources = append(requiredResources, r)
		}
	}

	// Validate resources have been applied.
	fmt.Println("Validate resources have been created")
	planfile, preNukeDiff, err := CreatePlan([]string{tfResourcePath})
	if err != nil {
		return err
	}
	if len(preNukeDiff) > 0 {
		fmt.Errorf("some resources have not been created. This is likely not caused by AWSNuke - %v", preNukeDiff)
	}

	// Destroy resources of given type.
	RunCmd("aws-nuke", "-c", "./example.yaml", "--force", "--profile", "default", "--no-dry-run", "--target", nukeType)

	// Remove resource from state to avoid ResourceNotFound error.
	RunCmd("terraform", "state", "rm", tfResourcePath)

	// Validate resources have been destroyed.
	planfile, postNukeDiff, err := CreatePlan([]string{tfResourcePath})
	if err != nil {
		return err
	}
	if len(postNukeDiff) < len(preCreateDiff)-len(requiredResources) {
		t.Errorf("not all resources created have been destroyed. The following still exist: %v", postNukeDiff)
	} else if len(postNukeDiff) > len(preCreateDiff) {
		t.Errorf("more resources have been destroyed than created. The AWS account was likely not empty")
	}
	return nil
}

func RunCmd(cmdName string, args ...string) error {
	if debugFlag {
		fmt.Println("os.exec:", cmdName, args)
	}
	cmd := exec.Command(cmdName, args...)
	if debugFlag {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreatePlan(targets []string) (*os.File, []string, error) {
	var diffResources []string

	f, err := ioutil.TempFile("/tmp", "AWSNuke")
	if err != nil {
		return nil, nil, err
	}

	params := []string{"plan", "-out=" + f.Name()}
	for _, t := range targets {
		params = append(params, "-target="+t)
	}

	if err := RunCmd("terraform", params...); err != nil {
		return nil, nil, err
	}

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return nil, nil, err
	}
	for _, m := range plan.Diff.Modules {
		for rName := range m.Resources {
			diffResources = append(diffResources, rName)
		}
	}
	return f, diffResources, nil
}

func TerraformApplyPlan(planfile *os.File) error {
	return RunCmd("terraform", "apply", "-input=false", "-auto-approve", planfile.Name())
}
