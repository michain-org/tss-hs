package msgHandle

import "fmt"

const (
	task = "Negotiation"
)

type (
	handleError struct {
		task string
		mode int
		err  error
	}
)

func (e *handleError) Error() string {
	return fmt.Sprintf("Mode:%d,Task:%s,Cause:%s", e.mode, e.task, e.err)
}

func newHandleError(mode int, task string, err error) *handleError {
	return &handleError{
		task: task,
		mode: mode,
		err:  err,
	}
}
