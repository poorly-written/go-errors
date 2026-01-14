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

type Code struct {
	http int
	grpc grpcCodes.Code
}

func (c Code) IsError() bool {
	return c.grpc != grpcCodes.OK
}

func (c Code) HttpCode() int {
	if c.http == 0 {
		return Ok.http
	}

	return c.http
}

func (c Code) GrpcCode() grpcCodes.Code {
	return c.grpc
}

var Ok = Code{http: 200, grpc: grpcCodes.OK}
var Created = Code{http: 201, grpc: grpcCodes.OK}
var Accepted = Code{http: 202, grpc: grpcCodes.OK}
var NoContent = Code{http: 204, grpc: grpcCodes.OK}
var PartialContent = Code{http: 206, grpc: grpcCodes.OK}

var BadRequest = Code{http: 400, grpc: grpcCodes.InvalidArgument}
var Unauthorized = Code{http: 401, grpc: grpcCodes.Unauthenticated}
var PaymentRequired = Code{http: 402, grpc: grpcCodes.FailedPrecondition}
var Forbidden = Code{http: 403, grpc: grpcCodes.PermissionDenied}
var NotFound = Code{http: 404, grpc: grpcCodes.NotFound}
var MethodNotAllowed = Code{http: 405, grpc: grpcCodes.Unimplemented}
var NotAcceptable = Code{http: 406, grpc: grpcCodes.InvalidArgument}
var ProxyAuthenticationRequired = Code{http: 407, grpc: grpcCodes.FailedPrecondition}
var RequestTimeout = Code{http: 408, grpc: grpcCodes.DeadlineExceeded}
var Conflict = Code{http: 409, grpc: grpcCodes.AlreadyExists}
var Gone = Code{http: 410, grpc: grpcCodes.NotFound}
var LengthRequired = Code{http: 411, grpc: grpcCodes.FailedPrecondition}
var PreconditionFailed = Code{http: 412, grpc: grpcCodes.FailedPrecondition}
var RequestEntityTooLarge = Code{http: 413, grpc: grpcCodes.FailedPrecondition}
var RequestUriTooLong = Code{http: 414, grpc: grpcCodes.FailedPrecondition}
var UnsupportedMediaType = Code{http: 415, grpc: grpcCodes.FailedPrecondition}
var RequestedRangeNotSatisfiable = Code{http: 416, grpc: grpcCodes.FailedPrecondition}
var ExpectationFailed = Code{http: 417, grpc: grpcCodes.FailedPrecondition}
var IAmATeapot = Code{http: 418, grpc: grpcCodes.FailedPrecondition}
var MisdirectedRequest = Code{http: 421, grpc: grpcCodes.FailedPrecondition}
var UnprocessableEntity = Code{http: 422, grpc: grpcCodes.InvalidArgument}
var Locked = Code{http: 423, grpc: grpcCodes.FailedPrecondition}
var FailedDependency = Code{http: 424, grpc: grpcCodes.FailedPrecondition}
var TooEarly = Code{http: 425, grpc: grpcCodes.FailedPrecondition}
var UpgradeRequired = Code{http: 426, grpc: grpcCodes.FailedPrecondition}
var PreconditionRequired = Code{http: 428, grpc: grpcCodes.FailedPrecondition}
var TooManyRequests = Code{http: 429, grpc: grpcCodes.ResourceExhausted}
var RequestHeaderFieldsTooLarge = Code{http: 431, grpc: grpcCodes.FailedPrecondition}
var UnavailableForLegalReasons = Code{http: 451, grpc: grpcCodes.FailedPrecondition}
var ClientClosedRequest = Code{http: 499, grpc: grpcCodes.Canceled}

var InternalServerError = Code{http: 500, grpc: grpcCodes.Internal}
var NotImplemented = Code{http: 501, grpc: grpcCodes.Unimplemented}
var BadGateway = Code{http: 502, grpc: grpcCodes.Unavailable}
var ServiceUnavailable = Code{http: 503, grpc: grpcCodes.Unavailable}
var GatewayTimeout = Code{http: 504, grpc: grpcCodes.DeadlineExceeded}
var VersionNotSupported = Code{http: 505, grpc: grpcCodes.Unimplemented}
var VariantAlsoNegotiatesExperimental = Code{http: 506, grpc: grpcCodes.FailedPrecondition}
var InsufficientStorage = Code{http: 507, grpc: grpcCodes.ResourceExhausted}
var LoopDetected = Code{http: 508, grpc: grpcCodes.Internal}
var NotExtended = Code{http: 510, grpc: grpcCodes.FailedPrecondition}
var NetworkAuthenticationRequired = Code{http: 511, grpc: grpcCodes.FailedPrecondition}
var WebServerReturnedAnUnknownError = Code{http: 520, grpc: grpcCodes.Unknown}
