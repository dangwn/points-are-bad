package main

import (
    "flag"
    "testing"
    vpc "points-are-bad/infra/vpc"
    tf_utils "points-are-bad/infra/tf_utils"
)

var (
    skipCleanup bool
    skipDestroy bool
)

func TestMain(t *testing.T) {
    flag.BoolVar(&skipCleanup, "skipCleanup", false, "Skip cleanup of local terraform files")
    flag.BoolVar(&skipDestroy, "skipDestroy", false, "Skip destroying infrastructure and leave it in AWS")

    flag.Parse()

    vpcStack := vpc.DeployVpc(t, skipCleanup, skipDestroy)
    defer vpcStack.TearDown(t)

    if !skipCleanup {
        tf_utils.RemoveTFCacheDir(".")
    }
}