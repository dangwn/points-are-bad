package main

import (
    "fmt"
    "flag"
    "testing"
    vpc "points-are-bad/infra/vpc"
    tf_utils "points-are-bad/infra/tf_utils"
)

var (
    skipDestroy bool
    skipCleanup bool
)

func TestMain(t *testing.T) {
    flag.BoolVar(&skipDestroy, "skipDestroy", false, "Skip destroying infrastructure and leave it in AWS")
    flag.BoolVar(&skipCleanup, "skipCleanup", false, "Skip cleanup of local terraform files")

    flag.Parse()

    outputs := vpc.DeployVpc(t, skipDestroy, skipCleanup)
    fmt.Println(fmt.Sprint("\n\n", outputs, "\n\n"))

    if !skipCleanup {
        tf_utils.RemoveTFCacheDir(".")
    }
}