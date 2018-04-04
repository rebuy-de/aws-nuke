package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform/terraform"
	"path"
	"strings"
)

type Terraform struct {
	// AWS region to use. E.g.: eu-central-1
	Region string
	// AWS profile to use. Must be defined in $HOME/.aws/credentials
	Profile string
	// Optional flag for printing the terraform command line output.
	Debug bool
	// ressDir is the absolute path of the directory where the Terraform resource
	// files reside.
	ressDir string
	// cmdDir is the path of the temporary directory terraform will be executed
	// from and where the meta files will reside.
	cmdDir string
}

// New creates a new terraform struct from given region, profile and relative directory.
// The terraform configuration and statefile will be created at
// /tmp/AWSNukeTest<random>.
func New(region, profile, dir string) (*Terraform, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cmdDir, err := ioutil.TempDir("/tmp", "AWSNukeTest")
	if err != nil {
		return nil, err
	}

	t := &Terraform{
		Region:  region,
		Profile: profile,
		ressDir: path.Join(wd, dir),
		cmdDir:  cmdDir,
	}
	return t, nil
}

func (t *Terraform) Init() error {
	return t.run(
		[]string{"init"},
		[]string{"-input=false", "-no-color"},
		[]string{t.ressDir},
	)
}

func (t *Terraform) PrintVersion() error {
	cmd := exec.Command("terraform", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *Terraform) CreatePlan(targets []string) (*os.File, []string, error) {
	var diffResources []string

	f, err := ioutil.TempFile(t.cmdDir, "AWSNuke")
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	// TODO(bethge): Should we eventually delete the tmp file?

	params := []string{"-out=" + f.Name()}
	for _, t := range targets {
		params = append(params, "-target="+t)
	}

	err = t.run(
		[]string{"plan"},
		params,
		[]string{t.ressDir},
	)
	if err != nil {
		return nil, nil, err
	}

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return nil, nil, err
	}
	for _, m := range plan.Diff.Modules {
		for name := range m.Resources {
			diffResources = append(diffResources, name)
		}
	}
	return f, diffResources, nil
}

func (t *Terraform) ApplyPlan(pFile *os.File) error {
	return t.run(
		[]string{"apply"},
		[]string{"-input=false", "-auto-approve", "-backup=-"},
		[]string{pFile.Name()},
	)
}

func (t *Terraform) RemoveResourceFromState(resPath string) error {
	return t.run(
		[]string{"state", "rm"},
		nil,
		[]string{resPath},
	)
}

// RemoveAllResources removes all terraform managed resources.
func (t *Terraform) RemoveAllResources() error {
	return t.run(
		[]string{"destroy"},
		[]string{"-input=false", "-no-color", "-auto-approve"},
		[]string{t.ressDir},
	)
}

func (t *Terraform) CreateResource(resPath string) error {
	planfile, preCreateDiff, err := t.CreatePlan([]string{resPath})
	if err != nil {
		return err
	}
	if len(preCreateDiff) == 0 {
		return fmt.Errorf(`no resources to create. The AWS resource may already exist or the terraform resource cannot be found: "%s"`, resPath)
	}

	// Terraform also creates the resource's dependencies.
	return t.ApplyPlan(planfile)
}

func (t *Terraform) RemoveResource(resPath string) error {
	return t.run(
		[]string{"destroy"},
		[]string{"-input=false", "-no-color", "-auto-approve", "-target=" + resPath},
		[]string{t.ressDir},
	)
}

// ResourceProperty retrieves the value for given property of the given
// terraform-managed AWS resource.
func (t *Terraform) ResourceProperty(resPath, prop string) (string, error) {
	c := t.cmd(
		// "terraform output" and piping into "terraform console" caused panic.
		// This is probably due to the configuration and state files being
		// in a different directory than the resources.
		[]string{"state", "show"},
		[]string{"-no-color"},
		[]string{resPath},
	)
	c.Stderr = os.Stderr
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	for _, r := range strings.Split(string(out), "\n") {
		kv := strings.SplitN(r, "=", 2)
		if len(kv) == 2 && prop == strings.TrimSpace(kv[0]) {
			return strings.TrimSpace(kv[1]), nil
		}
	}

	return "", fmt.Errorf("property '%v' not found for '%v'", prop, resPath)
}

func (t *Terraform) run(cmds, opts, args []string) error {
	if t.Debug {
		fmt.Println("os.exec:", "terraform", cmds, opts, args)
	}
	cmd := t.cmd(cmds, opts, args)
	if t.Debug {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *Terraform) cmd(cmds, opts, args []string) *exec.Cmd {
	args = append(append(cmds, opts...), args...)
	c := exec.Command("terraform", args...)
	c.Env = append(os.Environ(), []string{
		"TF_IN_AUTOMATION=1",
		"AWS_PROFILE=" + t.Profile,
		"AWS_DEFAULT_REGION=" + t.Region}...,
	)
	c.Dir = t.cmdDir
	return c
}
