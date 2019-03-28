package controller

import (
	"github.com/redhat-cop/events-notifier/pkg/controller/event"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, event.Add)
}
