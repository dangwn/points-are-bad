package iam

import (
	tf_utils "points-are-bad/infra/tf_utils"
	"testing"
)

func DeployIam(t *testing.T, skipCleanup, skipDestroy bool, region string) tf_utils.StackData {
	rootDir := "./iam"
	iamData := tf_utils.StackData{
		TerraformDir: rootDir,
		SkipCleanup: skipCleanup,
		SkipDestroy: skipDestroy,
		StackParams: map[string]interface{}{
			"aws_region": region,
			"eks_cluster_name": "points-are-bad-eks-cluster",
			"eks_cluster_role_name": "points-are-bad-eks-cluster-role",
			"eks_lambda_deploy_role_name": "points-are-bad-eks-lambda-deploy-role",
			"eks_master_role_name": "points-are-bad-eks-master-role",
		},
	}

	iamData.DeployStack(t)
	iamData.CollectOutputs(t)

	return iamData
}