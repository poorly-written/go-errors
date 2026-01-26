package errors

import (
	grpcCodes "google.golang.org/grpc/codes"
)

func codeFromValue(v int) Code {
	if v, ok := codeMapper[v]; ok {
		return v
	}

	return defaultErrorCodeIfNotInMapper
}

type Code interface {
	IsError() bool
	HttpCode() int
	GrpcCode() grpcCodes.Code
}

type code struct {
	http int
	grpc grpcCodes.Code
}

func (c code) IsError() bool {
	return c.grpc != grpcCodes.OK
}

func (c code) HttpCode() int {
	return c.http
}

func (c code) GrpcCode() grpcCodes.Code {
	return c.grpc
}

var Ok = code{http: 200, grpc: grpcCodes.OK}
var Created = code{http: 201, grpc: grpcCodes.OK}
var Accepted = code{http: 202, grpc: grpcCodes.OK}
var NoContent = code{http: 204, grpc: grpcCodes.OK}
var PartialContent = code{http: 206, grpc: grpcCodes.OK}

var BadRequest = code{http: 400, grpc: grpcCodes.InvalidArgument}
var Unauthorized = code{http: 401, grpc: grpcCodes.Unauthenticated}
var PaymentRequired = code{http: 402, grpc: grpcCodes.FailedPrecondition}
var Forbidden = code{http: 403, grpc: grpcCodes.PermissionDenied}
var NotFound = code{http: 404, grpc: grpcCodes.NotFound}
var MethodNotAllowed = code{http: 405, grpc: grpcCodes.Unimplemented}
var NotAcceptable = code{http: 406, grpc: grpcCodes.InvalidArgument}
var ProxyAuthenticationRequired = code{http: 407, grpc: grpcCodes.FailedPrecondition}
var RequestTimeout = code{http: 408, grpc: grpcCodes.DeadlineExceeded}
var Conflict = code{http: 409, grpc: grpcCodes.AlreadyExists}
var Gone = code{http: 410, grpc: grpcCodes.NotFound}
var LengthRequired = code{http: 411, grpc: grpcCodes.FailedPrecondition}
var PreconditionFailed = code{http: 412, grpc: grpcCodes.FailedPrecondition}
var RequestEntityTooLarge = code{http: 413, grpc: grpcCodes.FailedPrecondition}
var RequestUriTooLong = code{http: 414, grpc: grpcCodes.FailedPrecondition}
var UnsupportedMediaType = code{http: 415, grpc: grpcCodes.FailedPrecondition}
var RequestedRangeNotSatisfiable = code{http: 416, grpc: grpcCodes.FailedPrecondition}
var ExpectationFailed = code{http: 417, grpc: grpcCodes.FailedPrecondition}
var IAmATeapot = code{http: 418, grpc: grpcCodes.FailedPrecondition}
var MisdirectedRequest = code{http: 421, grpc: grpcCodes.FailedPrecondition}
var UnprocessableEntity = code{http: 422, grpc: grpcCodes.InvalidArgument}
var Locked = code{http: 423, grpc: grpcCodes.FailedPrecondition}
var FailedDependency = code{http: 424, grpc: grpcCodes.FailedPrecondition}
var TooEarly = code{http: 425, grpc: grpcCodes.FailedPrecondition}
var UpgradeRequired = code{http: 426, grpc: grpcCodes.FailedPrecondition}
var PreconditionRequired = code{http: 428, grpc: grpcCodes.FailedPrecondition}
var TooManyRequests = code{http: 429, grpc: grpcCodes.ResourceExhausted}
var RequestHeaderFieldsTooLarge = code{http: 431, grpc: grpcCodes.FailedPrecondition}
var UnavailableForLegalReasons = code{http: 451, grpc: grpcCodes.FailedPrecondition}
var ClientClosedRequest = code{http: 499, grpc: grpcCodes.Canceled}

var InternalServerError = code{http: 500, grpc: grpcCodes.Internal}
var NotImplemented = code{http: 501, grpc: grpcCodes.Unimplemented}
var BadGateway = code{http: 502, grpc: grpcCodes.Unavailable}
var ServiceUnavailable = code{http: 503, grpc: grpcCodes.Unavailable}
var GatewayTimeout = code{http: 504, grpc: grpcCodes.DeadlineExceeded}
var VersionNotSupported = code{http: 505, grpc: grpcCodes.Unimplemented}
var VariantAlsoNegotiatesExperimental = code{http: 506, grpc: grpcCodes.FailedPrecondition}
var InsufficientStorage = code{http: 507, grpc: grpcCodes.ResourceExhausted}
var LoopDetected = code{http: 508, grpc: grpcCodes.Internal}
var NotExtended = code{http: 510, grpc: grpcCodes.FailedPrecondition}
var NetworkAuthenticationRequired = code{http: 511, grpc: grpcCodes.FailedPrecondition}
var WebServerReturnedAnUnknownError = code{http: 520, grpc: grpcCodes.Unknown}
