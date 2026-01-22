package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"google.golang.org/grpc/metadata"
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
	StackFrames() []frame
	EnableShouldReport() DetailedError
	ShouldReport() bool
	Code(code Code) DetailedError
	InternalCode(errorCode string) DetailedError
	AddMetadata(key string, value interface{}) DetailedError
	GetMetadata() map[string]interface{}
	HasMetadata(keys ...string) bool
}

type err struct {
	message      string
	original     error
	frames       []frame
	stErr        *status.Status
	headers      metadata.MD
	trailers     metadata.MD
	reportable   bool
	code         Code
	internalCode *string
	metadata     map[string]interface{}
}

func (e *err) Error() string {
	return e.message
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

func (e *err) StackFrames() []frame {
	return e.frames
}

func (e *err) EnableShouldReport() DetailedError {
	e.reportable = true

	return e
}

func (e *err) ShouldReport() bool {
	return e.reportable
}

func (e *err) Code(code Code) DetailedError {
	e.code = code

	return e
}

func (e *err) InternalCode(code string) DetailedError {
	e.internalCode = &code

	return e
}

func (e *err) AddMetadata(key string, value interface{}) DetailedError {
	e.metadata[key] = value

	return e
}

func (e *err) GetMetadata() map[string]interface{} {
	return e.metadata
}

func (e *err) HasMetadata(keys ...string) bool {
	if len(keys) == 0 {
		return len(e.metadata) > 0
	}

	_, ok := e.metadata[keys[0]]

	return ok
}

func New(e interface{}, opts ...ErrorOption) DetailedError {
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
	case nil:
		message = ""
	default:
		message = fmt.Sprintf("%v", e)
	}

	errOpts := &errorOptions{
		headers:      make(metadata.MD),
		trailers:     make(metadata.MD),
		callerOffset: 2,
		message:      message,
	}

	for _, opt := range opts {
		opt.apply(errOpts)
	}

	callers := make([]uintptr, stackTraceDepth)
	length := runtime.Callers(errOpts.callerOffset, callers[:])

	frames := make([]frame, length)
	for i, pc := range callers[:length] {
		frames[i] = frame(pc)
	}

	de := &err{
		message:      errOpts.message,
		original:     original,
		frames:       frames,
		stErr:        nil,
		headers:      errOpts.headers,
		trailers:     errOpts.trailers,
		code:         defaultErrorCode,
		reportable:   false,
		internalCode: nil,
		metadata:     make(map[string]interface{}),
	}

	stErr, ok := status.FromError(original)
	if !ok || stErr == nil {
		return de
	}

	statusCode := int(stErr.Code())
	if httpStatusCode := errOpts.headers.Get(httpHeaderKey); len(httpStatusCode) > 0 {
		if cc, e := strconv.Atoi(httpStatusCode[0]); e == nil {
			statusCode = cc
		}
	}

	// if a message is not provided and message from error is not an empty string, overwrite it
	if stMsg := stErr.Message(); message == "" && stMsg != "" {
		de.message = stMsg
	}

	if code := codeFromValue(statusCode); code.IsError() {
		de.code = code
	}

	return de
}
