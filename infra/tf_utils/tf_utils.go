package tf_utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

type StackData struct {
    TerraformDir    string
    StackParams     map[string]interface{}
    SkipCleanup     bool
    SkipDestroy     bool

    StackOutputs    map[string]interface{}
    
    Options         *terraform.Options
}

func (sd *StackData) configureOptions() {
    if sd.Options == nil {
        sd.Options = &terraform.Options{
            TerraformDir:   sd.TerraformDir,
            Vars:           sd.StackParams,
        }
    }
}

func (sd StackData) DeployStack(t *testing.T) {
    sd.configureOptions()
    terraform.InitAndApply(t, sd.Options)
}

func (sd *StackData) CollectOutputs(t *testing.T) {
    sd.configureOptions()
    outputs := terraform.OutputAll(t, sd.Options)
    sd.StackOutputs = outputs
}

func (sd StackData) DestroyStack(t *testing.T) {
    sd.configureOptions()
    terraform.Destroy(t, sd.Options)
}

func (sd StackData) TearDown(t *testing.T) {
    if !sd.SkipDestroy {
        sd.DestroyStack(t)
    }
    if !sd.SkipCleanup {
        RemoveTerraformStaticFiles(sd.TerraformDir)
        RemoveTFCacheDir(sd.TerraformDir)
    }
}

func RemoveTFCacheDir(rootDir string) {
    tfCachePath := rootDir + "/.terraform"
    if err := os.RemoveAll(tfCachePath); err != nil {
        fmt.Println(err)
    }
}

func RemoveTerraformStaticFiles(rootDir string) {
    var tfStateFiles [3]string
	tfStateFiles[0] = rootDir + "/terraform.tfstate.backup"
	tfStateFiles[1] = rootDir + "/terraform.tfstate"
	tfStateFiles[2] = rootDir + "/.terraform.lock.hcl"
    for _, f := range tfStateFiles {
        if err := os.Remove("./" + f); err != nil {
            fmt.Println(err)
        }
    }
}