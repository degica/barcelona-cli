package operations

import (
	"github.com/urfave/cli"
)

type Operation interface {
	run() *runResult
}

type runResult struct {
	is_error bool
	message  string
}

func Execute(oper Operation) error {
	run_result := oper.run()

	if run_result == nil {
		return nil
	}

	if run_result.is_error {
		return cli.NewExitError(run_result.message, 1)
	}

	return nil
}

func error_result(message string) *runResult {
	return &runResult{
		message:  message,
		is_error: true,
	}
}

func ok_result() *runResult {
	return &runResult{
		message:  "",
		is_error: false,
	}
}

type OperationType string

const (
	Delete OperationType = "Delete"
	Show                 = "Show"
)
