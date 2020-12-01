package operations

import (
	"fmt"
)

type ProfileOperation struct {
	opname string
	name string
}

func NewProfileOperation(opname string, name string) *ProfileOperation {
	return &ProfileOperation{
		opname: opname,
		name: name,
	}
}

func (oper ProfileOperation) run() *runResult {
	switch oper.opname {
	case "create":
		return createProfile(oper.name)
	case "delete":
		return deleteProfile(oper.name)
	case "use":
		return useProfile(oper.name)
	case "show":
		return showProfile(oper.name)
	}
	return error_result("Unknown command")
}

func createProfile(name string) *runResult {
	if name == "" {
		return error_result("Please enter a name")
	}
	
	return ok_result()
}

func deleteProfile(name string) *runResult {
	if name == "" {
		return error_result("Please enter a name")
	}
	return ok_result()
}

func useProfile(name string) *runResult {
	if name == "" {
		return error_result("Please enter a name")
	}
	return ok_result()
}

func showProfile(name string) *runResult {
	if name == "" {
		fmt.Println("Current profile: default")
	}
	return ok_result()
}
