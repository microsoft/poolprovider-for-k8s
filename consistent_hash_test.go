package main

import (
	"testing"
)

func TestConsistentHashShouldReturnCorrectValueBasedOnKey(t *testing.T) {
	nodes := []string{"buildkitd-0", "buildkitd-1"}

	chosen:= ComputeConsistentHash(nodes, "microsoft/k8spoolprovider.")
	
	chosen2:= ComputeConsistentHash(nodes, "microsoft/k8spoolproviderroot/src.")
	if (chosen != "buildkitd-1" || chosen2 != "buildkitd-0"){
		t.Errorf("Consistent Hashing failed")
	}

}

func TestConsistentHashShouldReturnSameValueBasedOnKey(t *testing.T) {
	nodes := []string{"buildkitd-0", "buildkitd-1"}

	chosen:= ComputeConsistentHash(nodes, "microsoft/k8spoolprovider.")
	
	chosen2:= ComputeConsistentHash(nodes, "microsoft/k8spoolprovider.")
	if (chosen != "buildkitd-1" || chosen2 != "buildkitd-1"){
		t.Errorf("Consistent Hashing failed")
	}

}