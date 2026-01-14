package errors

import (
	"errors"
	"fmt"
	"runtime"

	"google.golang.org/grpc/status"
)

type DetailedError interface {
	error
	Unwrap() error
	Original() error
	GRPCStatus() *status.Status
	HasError() bool
	AddMessage(msg string) DetailedError
	AddHeader(key string, value ...string) DetailedError
	RemoveHeader(key string) DetailedError
	AddTrailer(key string, value ...string) DetailedError
	RemoveTrailer(key string) DetailedError
	StackFrames() []uintptr
}

type err struct {
	message  string
	original error
	frames   []uintptr
	stErr    *status.Status
	headers  map[string][]string
	trailers map[string][]string
}

func (e *err) Error() string {
	if e.original == nil {
		return e.message
	}

	return e.original.Error()
}

func (e *err) Unwrap() error {
	if e.original != nil {
		return e.original
	}

	return nil
}

func (e *err) Original() error {
	return e.original
}

func (e *err) GRPCStatus() *status.Status {
	return e.stErr
}

func (e *err) HasError() bool {
	return e.original != nil
}

func (e *err) AddMessage(msg string) DetailedError {
	e.message = msg

	return e
}

func (e *err) AddHeader(key string, value ...string) DetailedError {
	e.headers[key] = append(e.headers[key], value...)

	return e
}

func (e *err) RemoveHeader(key string) DetailedError {
	delete(e.headers, key)

	return e
}

func (e *err) AddTrailer(key string, value ...string) DetailedError {
	e.trailers[key] = append(e.trailers[key], value...)

	return e
}

func (e *err) RemoveTrailer(key string) DetailedError {
	delete(e.trailers, key)

	return e
}

func (e *err) StackFrames() []uintptr {
	return e.frames
}

func New(e interface{}, msg ...string) DetailedError {
	var original error
	var message string
	switch e := e.(type) {
	case error:
		// Firstly, if it's another DetailedError instance, then return early
		// otherwise process it
		var de DetailedError
		if errors.As(e, &de) {
			return de
		}

		original = e
		message = e.Error()
	default:
		message = fmt.Sprintf("%v", e)
	}

	if len(msg) > 0 {
		message = msg[0]
	}

	callers := make([]uintptr, 50)
	length := runtime.Callers(2, callers[:])

	de := &err{
		message:  message,
		original: original,
		frames:   callers[:length],
		headers:  make(map[string][]string),
		trailers: make(map[string][]string),
	}

	stErr, ok := status.FromError(original)
	if !ok || stErr == nil {
		return de
	}

	de.stErr = stErr

	return de
}
