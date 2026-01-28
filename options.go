package errors

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type errorOptions struct {
	Message      string
	Headers      metadata.MD
	Trailers     metadata.MD
	CallerOffset int
	Ctx          context.Context
	InternalCode *string
	Code         Code
	Reportable   bool
	SkipOnNil    bool
}

type ErrorOption interface {
	apply(error, *errorOptions)
}

type errorOptionModifier func(error, *errorOptions)

type funcErrorOption struct {
	f errorOptionModifier
}

func (fo *funcErrorOption) apply(err error, opt *errorOptions) {
	fo.f(err, opt)
}

func newFuncErrorOption(f errorOptionModifier) *funcErrorOption {
	return &funcErrorOption{
		f: f,
	}
}

func Message(msg string) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.Message = msg
	})
}

func Headers(headers metadata.MD) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.Headers = headers
	})
}

func Trailers(trailers metadata.MD) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.Trailers = trailers
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
			o.CallerOffset = callerOffset
		} else {
			o.CallerOffset += callerOffset
		}
	})
}

func Context(ctx context.Context) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.Ctx = ctx
	})
}

func InternalCode(code string) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.InternalCode = &code
	})
}

func ErrorCode(code Code) ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.Code = code
	})
}

func Reportable() ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.Reportable = true
	})
}

func SkipOnNil() ErrorOption {
	return newFuncErrorOption(func(_ error, o *errorOptions) {
		o.SkipOnNil = true
	})
}

func ErrorOptionModifier(modifier errorOptionModifier) ErrorOption {
	return newFuncErrorOption(func(err error, o *errorOptions) {
		modifier(err, o)
	})
}
