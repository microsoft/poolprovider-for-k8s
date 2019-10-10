package main

import (
	"github.com/serialx/hashring"
)

func ComputeConsistentHash(nodes []string, key string) (string) {
	ring := hashring.New(nodes)
	x, ok := ring.GetNode(key)
	if !ok {
		return ""
	}
	return x
}