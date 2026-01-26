package errors

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type DetailedError interface {
	error
	Unwrap() error
	Original() error
	GRPCStatus() *status.Status
	HasError() bool
	Message(msg string) DetailedError
	AddHeader(key string, value ...string) DetailedError
	RemoveHeader(key string) DetailedError
	AddTrailer(key string, value ...string) DetailedError
	RemoveTrailer(key string) DetailedError
	StackFrames() []frame
	ShouldBeReported() DetailedError
	IsReportable() bool
	Code(code Code) DetailedError
	InternalCode(errorCode string) DetailedError
	Context(ctx context.Context, extractMetadata ...bool) DetailedError
	AddMetadata(key string, value interface{}) DetailedError
	GetMetadata() map[string]interface{}
	HasMetadata(keys ...string) bool
	IncludeMetadata() DetailedError
	AddReason(key string, reason Reason) DetailedError
	GetReasons() map[string][]Reason
	HasReasons(keys ...string) bool
	Append(key string, value interface{}) DetailedError
	Send()
}

type err struct {
	status          *status.Status
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

func (e *err) GRPCStatus() *status.Status {
	st := status.New(e.code.GrpcCode(), e.message)

	marshaled, err := errorMarshaler(&ErrorDetails{
		Message:         &e.message,
		InternalCode:    e.internalCode,
		Reasons:         e.reasons,
		IncludeMetadata: e.includeMetadata,
		Metadata:        e.metadata,
	})

	// error occurred during error marshalling
	if err != nil {
		return status.New(grpcCodes.Internal, err.Error())
	}

	if marshaled == nil {
		e.status = st

		return st
	}

	dSt, err := st.WithDetails(marshaled)
	if err != nil {
		return status.New(grpcCodes.Internal, err.Error())
	}

	e.status = dSt

	return dSt
}

func (e *err) HasError() bool {
	return e.original != nil
}

func (e *err) Message(msg string) DetailedError {
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

func (e *err) ShouldBeReported() DetailedError {
	e.reportable = true

	return e
}

func (e *err) IsReportable() bool {
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

func (e *err) Context(ctx context.Context, extractMetadata ...bool) DetailedError {
	e.ctx = ctx

	if len(extractMetadata) == 0 || extractMetadata[0] == false {
		return e
	}

	for k, v := range contextualMetadataExtractor(ctx) {
		e.AddMetadata(k, v)
	}

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
	if len(e.metadata) == 0 {
		return false
	}

	for _, key := range keys {
		if _, ok := e.metadata[key]; !ok {
			return false
		}
	}

	return true
}

func (e *err) AddReason(key string, reason Reason) DetailedError {
	if _, ok := e.reasons[key]; !ok {
		e.reasons[key] = make([]Reason, 0)
	}

	e.reasons[key] = append(e.reasons[key], reason)

	return e
}

func (e *err) HasReasons(keys ...string) bool {
	if len(e.reasons) == 0 {
		return false
	}

	for _, key := range keys {
		if _, ok := e.reasons[key]; !ok {
			return false
		}
	}

	return true
}

func (e *err) GetReasons() map[string][]Reason {
	return e.reasons
}

// Append method appends either to reasons or metadata based on the value provided.
// If the value is a type of `Reason`, then append forwards the call to the
// `AddReason` function. Otherwise, it forwards the call to the `AddMetadata` function.
func (e *err) Append(key string, value interface{}) DetailedError {
	if r, ok := value.(Reason); ok {
		return e.AddReason(key, r)
	}

	return e.AddMetadata(key, value)
}

func (e *err) Send() {
	// set http status code header
	e.headers.Set(httpHeaderKey, strconv.Itoa(e.code.HttpCode()))

	if e.headers.Len() > 0 {
		grpc.SetHeader(e.ctx, e.headers)
	}

	if e.trailers.Len() > 0 {
		grpc.SetTrailer(e.ctx, e.trailers)
	}
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
		message:      message,
		headers:      make(metadata.MD),
		trailers:     make(metadata.MD),
		callerOffset: 2,
		ctx:          context.Background(),
		internalCode: nil,
		code:         defaultErrorCode,
		reportable:   false,
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
		code:         errOpts.code,
		reportable:   errOpts.reportable,
		internalCode: errOpts.internalCode,
		metadata:     make(map[string]interface{}),
		ctx:          errOpts.ctx,
	}

	stErr, ok := status.FromError(original)
	if !ok || stErr == nil {
		return de
	}

	statusCode := int(stErr.Code())
	if httpStatusCode := errOpts.headers.Get(httpHeaderKey); len(httpStatusCode) > 0 {
		// as header key is present, we're removing it.
		// it will be added later in `GRPCStatus` method because of the status code
		de.headers.Delete(httpHeaderKey)

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
