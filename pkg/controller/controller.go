package controller

import (
	eventv1 "github.com/redhat-cop/events-notifier/pkg/apis/event/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager, *[]eventv1.EventSubscription) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, es *[]eventv1.EventSubscription) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m, es); err != nil {
			return err
		}
	}
	return nil
}
