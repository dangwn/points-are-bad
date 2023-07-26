package kms

import (
	tf_utils "points-are-bad/infra/tf_utils"
	"testing"
)

func DeployKms(t *testing.T, skipCleanup, skipDestroy bool, region string) tf_utils.StackData {
	rootDir := "./kms"
	kmsData := tf_utils.StackData{
		TerraformDir:	rootDir,
		SkipCleanup:	skipCleanup,
		SkipDestroy:	skipDestroy,
		StackParams:	map[string]interface{}{
			"aws_region":	region,
		},
	}

	kmsData.DeployStack(t)
	kmsData.CollectOutputs(t)

	return kmsData
}