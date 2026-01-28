package errors

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type errorOptions struct {
	message      string
	headers      metadata.MD
	trailers     metadata.MD
	callerOffset int
	ctx          context.Context
	internalCode *string
	code         Code
	reportable   bool
	skipOnNil    bool
}

type ErrorOption interface {
	apply(error, *errorOptions)
}

type funcErrorOption struct {
	f func(error, *errorOptions)
}

func (fo *funcErrorOption) apply(err error, opt *errorOptions) {
	fo.f(err, opt)
}

func newFuncErrorOption(f func(error, *errorOptions)) *funcErrorOption {
	return &funcErrorOption{
		f: f,
	}
}

func Message(msg string) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.message = msg
	})
}

func Headers(headers metadata.MD) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.headers = headers
	})
}

func Trailers(trailers metadata.MD) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.trailers = trailers
	})
}

func CallerOffset(callerOffset int, startsFromOffset ...bool) ErrorOption {
	fromOffset := false
	if len(startsFromOffset) > 0 && startsFromOffset[0] {
		fromOffset = true
	}

	return newFuncErrorOption(func(_ error, o *errorOptions) {
		if callerOffset < 0 {
			callerOffset = 0
		}

		if fromOffset {
			o.callerOffset = callerOffset
		} else {
			o.callerOffset += callerOffset
		}
	})
}

func Context(ctx context.Context) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.ctx = ctx
	})
}

func InternalCode(code string) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.internalCode = &code
	})
}

func ErrorCode(code Code) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.code = code
	})
}

func Reportable() ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.reportable = true
	})
}

func SkipOnNil() ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.skipOnNil = true
	})
}
