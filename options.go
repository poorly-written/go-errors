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
}

type ErrorOption interface {
	apply(*errorOptions)
}

type funcErrorOption struct {
	f func(*errorOptions)
}

func (fo *funcErrorOption) apply(opt *errorOptions) {
	fo.f(opt)
}

func newFuncErrorOption(f func(*errorOptions)) *funcErrorOption {
	return &funcErrorOption{
		f: f,
	}
}

func Message(msg string) ErrorOption {
	return newFuncErrorOption(func(o *errorOptions) {
		o.message = msg
	})
}

func Headers(headers metadata.MD) ErrorOption {
	return newFuncErrorOption(func(o *errorOptions) {
		o.headers = headers
	})
}

func Trailers(trailers metadata.MD) ErrorOption {
	return newFuncErrorOption(func(o *errorOptions) {
		o.trailers = trailers
	})
}

func CallerOffset(callerOffset int, startsFromOffset ...bool) ErrorOption {
	fromOffset := false
	if len(startsFromOffset) > 0 && startsFromOffset[0] {
		fromOffset = true
	}

	return newFuncErrorOption(func(o *errorOptions) {
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
	return newFuncErrorOption(func(o *errorOptions) {
		o.ctx = ctx
	})
}

func InternalCode(code string) ErrorOption {
	return newFuncErrorOption(func(o *errorOptions) {
		o.internalCode = &code
	})
}

func ErrorCode(code Code) ErrorOption {
	return newFuncErrorOption(func(o *errorOptions) {
		o.code = code
	})
}
