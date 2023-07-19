package vpc

import (
	tf_utils "points-are-bad/infra/tf_utils"
	"testing"
)

func DeployVpc(t *testing.T, skipCleanup, skipDestroy bool) tf_utils.StackData {
    rootDir := "./vpc"
    vpcData := tf_utils.StackData{
        TerraformDir: rootDir,
        SkipCleanup: skipCleanup,
        SkipDestroy: skipDestroy,
        StackParams: map[string]interface{}{
            "aws_region": "us-east-1",
            "aws_az_1": "us-east-1a",
            "aws_az_2": "us-east-1b",
            "aws_az_3": "us-east-1c",
            "vpc_cidr_range": "10.0.0.0/16",
            "private_subnet_1_cidr_range": "10.0.0.0/24",
            "private_subnet_2_cidr_range": "10.0.1.0/24",
            "private_subnet_3_cidr_range": "10.0.2.0/24",
            "eks_endpoint_subnet_1_cidr_range": "10.0.3.0/24",
            "eks_endpoint_subnet_2_cidr_range": "10.0.4.0/24",
            "eks_endpoint_subnet_3_cidr_range": "10.0.5.0/24",
            "public_subnet_1_cidr_range": "10.0.6.0/24",
            "public_subnet_2_cidr_range": "10.0.7.0/24",
            "public_subnet_3_cidr_range": "10.0.8.0/24",
        },
    }

    vpcData.DeployStack(t)
    vpcData.CollectOutputs(t)

    return vpcData
}