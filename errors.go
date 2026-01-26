package errors

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"google.golang.org/grpc"
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
	Context(ctx context.Context) DetailedError
	IncludeMetadata() DetailedError
	AddMetadata(key string, value interface{}) DetailedError
	GetMetadata() map[string]interface{}
	HasMetadata(keys ...string) bool
	AddReason(key string, reason Reason) DetailedError
	HasReasons() bool
	GetReasons() map[string][]Reason
}

type err struct {
	message         string
	original        error
	frames          []frame
	headers         metadata.MD
	trailers        metadata.MD
	reasons         map[string][]Reason
	reportable      bool
	code            Code
	internalCode    *string
	metadata        map[string]interface{}
	includeMetadata bool
	ctx             context.Context
}

func (e *err) Error() string {
	return e.message
}

func (e *err) Unwrap() error {
	return e.original
}

func (e *err) Original() error {
	return e.original
}

func (e *err) setHeaders(md metadata.MD) error {
	return grpc.SetHeader(e.ctx, md)
}

func (e *err) setTrailers(md metadata.MD) error {
	return grpc.SetTrailer(e.ctx, md)
}

func (e *err) GRPCStatus() *status.Status {
	// set http status code header
	e.setHeaders(metadata.Pairs(httpHeaderKey, strconv.Itoa(e.code.HttpCode())))

	if e.headers.Len() > 0 {
		e.setHeaders(e.headers)
	}

	if e.trailers.Len() > 0 {
		e.setTrailers(e.headers)
	}

	st := status.New(e.code.GrpcCode(), e.message)

	marshaled, err := errorMarshaler(&ErrorDetails{
		Message:         &e.message,
		InternalCode:    e.internalCode,
		Reasons:         e.reasons,
		IncludeMetadata: e.includeMetadata,
		Metadata:        e.metadata,
	})
	if err != nil && marshaled != nil {
		st.WithDetails(marshaled)
	}

	return st
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

func (e *err) Context(ctx context.Context) DetailedError {
	e.ctx = ctx

	return e
}

func (e *err) IncludeMetadata() DetailedError {
	e.includeMetadata = true

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

func (e *err) AddReason(key string, reason Reason) DetailedError {
	if _, ok := e.reasons[key]; !ok {
		e.reasons[key] = make([]Reason, 0)
	}

	e.reasons[key] = append(e.reasons[key], reason)

	return e
}

func (e *err) HasReasons() bool {
	return len(e.reasons) > 0
}

func (e *err) GetReasons() map[string][]Reason {
	return e.reasons
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
		ctx:          nil,
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
		headers:      errOpts.headers,
		trailers:     errOpts.trailers,
		reasons:      make(map[string][]Reason),
		code:         defaultErrorCode,
		reportable:   false,
		internalCode: nil,
		metadata:     make(map[string]interface{}),
		ctx:          errOpts.ctx,
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

	for idx, detail := range stErr.Details() {
		details, err := errorUnmarshaler(idx, detail)
		if err != nil || details == nil {
			continue
		}

		if details.Message != nil {
			de.message = *details.Message
		}

		if details.InternalCode != nil {
			de.internalCode = details.InternalCode
		}

		for k, v := range details.Reasons {
			de.reasons[k] = append(de.reasons[k], v...)
		}

		// metadata is always overwritten here, at least now, there is no intention to merge multiple metadata from details
		for k, v := range details.Metadata {
			de.metadata[k] = v
		}
	}

	return de
}
