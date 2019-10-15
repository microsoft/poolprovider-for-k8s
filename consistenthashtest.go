package main

import (
	"testing"
    "fmt"
)

func TestConsistentHash(t *testing.T) {
	nodes := []string{"buildkitd-0", "buildkitd-1"}

	chosen:= ComputeConsistentHash(nodes, "microsoft/k8spoolprovider.")
	fmt.Println("Selected Node", chosen)
	chosen2:= ComputeConsistentHash(nodes, "microsoft/k8spoolproviderroot/src.")
	fmt.Println("Selected Node", chosen2)
}