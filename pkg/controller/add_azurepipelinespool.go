package controller

import (
	"github.com/microsoft/k8s-poolprovider/pkg/controller/azurepipelinespool"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, azurepipelinespool.Add)
}
