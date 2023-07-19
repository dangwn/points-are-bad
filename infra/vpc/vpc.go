package vpc

import (
	tf_utils "points-are-bad/infra/tf_utils"
	"testing"
)

func DeployVpc(t *testing.T, skipCleanup, skipDestroy bool, region, az1, az2, az3 string) tf_utils.StackData {
    rootDir := "./vpc"
    vpcData := tf_utils.StackData{
        TerraformDir: rootDir,
        SkipCleanup: skipCleanup,
        SkipDestroy: skipDestroy,
        StackParams: map[string]interface{}{
            "aws_region": region,
            "aws_az_1": az1,
            "aws_az_2": az2,
            "aws_az_3": az3,
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