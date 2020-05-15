package controller

import (
	"github.com/wu0407/daemonset-operator/pkg/controller/xdaemonset"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, xdaemonset.Add)
}
