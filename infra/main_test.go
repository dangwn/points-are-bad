package main

import (
    "flag"
    "testing"
    "points-are-bad/infra/iam"
    "points-are-bad/infra/vpc"
    tf_utils "points-are-bad/infra/tf_utils"
)

const (
    AWS_REGION string = "us-east-1"
    AWS_AZ_1 string = "us-east-1a"
    AWS_AZ_2 string = "us-east-1b"
    AWS_AZ_3 string = "us-east-1b"
)

var (
    skipCleanup bool
    skipDestroy bool
)

func TestMain(t *testing.T) {
    flag.BoolVar(&skipCleanup, "skipCleanup", false, "Skip cleanup of local terraform files")
    flag.BoolVar(&skipDestroy, "skipDestroy", false, "Skip destroying infrastructure and leave it in AWS")

    flag.Parse()

    vpcStack := vpc.DeployVpc(t, skipCleanup, skipDestroy, AWS_REGION, AWS_AZ_1, AWS_AZ_2, AWS_AZ_3)
    defer vpcStack.TearDown(t)

    iamStack := iam.DeployIam(t, skipCleanup, skipDestroy, AWS_REGION)
    defer iamStack.TearDown(t)

    if !skipCleanup {
        tf_utils.RemoveTFCacheDir(".")
    }
}