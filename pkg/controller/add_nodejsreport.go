package controller

import (
	"github.com/sam-github/operator-nodejs/pkg/controller/nodejsreport"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, nodejsreport.Add)
}
