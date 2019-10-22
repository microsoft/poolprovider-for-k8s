package main

import (
	"stathat.com/c/consistent"
)

func ComputeConsistentHash(nodes []string, key string) (string) {
	c := consistent.New()

	for _, items := range nodes {
		s := items
		if s != "" {
			c.Add(s)
		}
	}

	x, err := c.Get(key)
	if err != nil {
		return ""
	}
	return x
}