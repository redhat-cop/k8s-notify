package controller

import (
	"github.com/redhat-cop/events-notifier/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager, *util.SharedResources) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, sr *util.SharedResources) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m, sr); err != nil {
			return err
		}
	}
	return nil
}
