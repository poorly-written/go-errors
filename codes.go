package errors

import (
	"errors"
	"sync"

	grpcCodes "google.golang.org/grpc/codes"
)

var codeMapSetOnce sync.Once
var httpHeaderSetOnce sync.Once
var defaultCodeSetOnce sync.Once

func SetCodeMapper(mapper map[int]Code) {
	codeMapSetOnce.Do(func() {
		codeMapper = mapper
	})
}

func SetHttpHeaderKey(key string) {
	httpHeaderSetOnce.Do(func() {
		httpHeaderKey = key
	})
}

func SetDefaultCode(code Code) {
	defaultCodeSetOnce.Do(func() {
		defaultCodeIfNotMatched = code
	})
}

var httpHeaderKey = "x-http-status"

var defaultCodeIfNotMatched = BadRequest

var codeMapper = map[int]Code{
	int(grpcCodes.OK):                 Ok,
	int(grpcCodes.Canceled):           ClientClosedRequest,
	int(grpcCodes.Unknown):            BadRequest,
	int(grpcCodes.InvalidArgument):    UnprocessableEntity,
	int(grpcCodes.DeadlineExceeded):   GatewayTimeout,
	int(grpcCodes.NotFound):           NotFound,
	int(grpcCodes.AlreadyExists):      BadRequest,
	int(grpcCodes.PermissionDenied):   Forbidden,
	int(grpcCodes.ResourceExhausted):  InternalServerError,
	int(grpcCodes.FailedPrecondition): UnprocessableEntity,
	int(grpcCodes.Aborted):            ClientClosedRequest,
	int(grpcCodes.OutOfRange):         InternalServerError,
	int(grpcCodes.Unimplemented):      InternalServerError,
	int(grpcCodes.Internal):           InternalServerError,
	int(grpcCodes.Unavailable):        ServiceUnavailable,
	int(grpcCodes.DataLoss):           InternalServerError,
	int(grpcCodes.Unauthenticated):    Unauthorized,
	200:                               Ok,
	201:                               Created,
	202:                               Accepted,
	204:                               NoContent,
	206:                               PartialContent,
	400:                               BadRequest,
	401:                               Unauthorized,
	402:                               PaymentRequired,
	403:                               Forbidden,
	404:                               NotFound,
	405:                               MethodNotAllowed,
	406:                               NotAcceptable,
	407:                               ProxyAuthenticationRequired,
	408:                               RequestTimeout,
	409:                               Conflict,
	410:                               Gone,
	411:                               LengthRequired,
	412:                               PreconditionFailed,
	413:                               RequestEntityTooLarge,
	414:                               RequestUriTooLong,
	415:                               UnsupportedMediaType,
	416:                               RequestedRangeNotSatisfiable,
	417:                               ExpectationFailed,
	418:                               IAmATeapot,
	421:                               MisdirectedRequest,
	422:                               UnprocessableEntity,
	423:                               Locked,
	424:                               FailedDependency,
	425:                               TooEarly,
	426:                               UpgradeRequired,
	428:                               PreconditionRequired,
	429:                               TooManyRequests,
	431:                               RequestHeaderFieldsTooLarge,
	451:                               UnavailableForLegalReasons,
	500:                               InternalServerError,
	501:                               NotImplemented,
	502:                               BadGateway,
	503:                               ServiceUnavailable,
	505:                               VersionNotSupported,
	506:                               VariantAlsoNegotiatesExperimental,
	507:                               InsufficientStorage,
	508:                               LoopDetected,
	510:                               NotExtended,
	511:                               NetworkAuthenticationRequired,
	520:                               WebServerReturnedAnUnknownError,
}

func codeFromValue(v int) (Code, error) {
	if v, ok := codeMapper[v]; ok {
		return v, nil
	}

	return defaultCodeIfNotMatched, errors.New("no map for the given status code")
}

type Code struct {
	http int
	grpc grpcCodes.Code
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

var BadRequest = Code{http: 400, grpc: grpcCodes.InvalidArgument}
var Unauthorized = Code{http: 401, grpc: grpcCodes.Unauthenticated}
var Forbidden = Code{http: 403, grpc: grpcCodes.PermissionDenied}
var NotFound = Code{http: 404, grpc: grpcCodes.NotFound}
var Conflict = Code{http: 409, grpc: grpcCodes.AlreadyExists}
var IAmATeapot = Code{http: 418, grpc: grpcCodes.FailedPrecondition}
var UnprocessableEntity = Code{http: 422, grpc: grpcCodes.InvalidArgument}
var Locked = Code{http: 423, grpc: grpcCodes.FailedPrecondition}
var FailedDependency = Code{http: 424, grpc: grpcCodes.FailedPrecondition}
var UpgradeRequired = Code{http: 426, grpc: grpcCodes.FailedPrecondition}
var TooManyRequests = Code{http: 429, grpc: grpcCodes.ResourceExhausted}
var ClientClosedRequest = Code{http: 499, grpc: grpcCodes.Canceled}

var InternalServerError = Code{http: 500, grpc: grpcCodes.Internal}
var ServiceUnavailable = Code{http: 503, grpc: grpcCodes.Unavailable}
var GatewayTimeout = Code{http: 504, grpc: grpcCodes.DeadlineExceeded}
var InsufficientStorage = Code{http: 507, grpc: grpcCodes.ResourceExhausted}
var WebServerReturnedAnUnknownError = Code{http: 520, grpc: grpcCodes.Unknown}

var PartialContent = Code{http: 206, grpc: grpcCodes.OK}
var PaymentRequired = Code{http: 402, grpc: grpcCodes.FailedPrecondition}
var MethodNotAllowed = Code{http: 405, grpc: grpcCodes.Unimplemented}
var NotAcceptable = Code{http: 406, grpc: grpcCodes.InvalidArgument}
var RequestTimeout = Code{http: 408, grpc: grpcCodes.DeadlineExceeded}
var ProxyAuthenticationRequired = Code{http: 407, grpc: grpcCodes.FailedPrecondition}
var Gone = Code{http: 410, grpc: grpcCodes.NotFound}
var LengthRequired = Code{http: 411, grpc: grpcCodes.FailedPrecondition}
var PreconditionFailed = Code{http: 412, grpc: grpcCodes.FailedPrecondition}
var RequestEntityTooLarge = Code{http: 413, grpc: grpcCodes.FailedPrecondition}
var RequestUriTooLong = Code{http: 414, grpc: grpcCodes.FailedPrecondition}
var UnsupportedMediaType = Code{http: 415, grpc: grpcCodes.FailedPrecondition}
var RequestedRangeNotSatisfiable = Code{http: 416, grpc: grpcCodes.FailedPrecondition}
var ExpectationFailed = Code{http: 417, grpc: grpcCodes.FailedPrecondition}
var MisdirectedRequest = Code{http: 421, grpc: grpcCodes.FailedPrecondition}
var TooEarly = Code{http: 425, grpc: grpcCodes.FailedPrecondition}
var PreconditionRequired = Code{http: 428, grpc: grpcCodes.FailedPrecondition}
var RequestHeaderFieldsTooLarge = Code{http: 431, grpc: grpcCodes.FailedPrecondition}
var UnavailableForLegalReasons = Code{http: 451, grpc: grpcCodes.FailedPrecondition}
var NotImplemented = Code{http: 501, grpc: grpcCodes.Unimplemented}
var BadGateway = Code{http: 502, grpc: grpcCodes.Unavailable}
var VersionNotSupported = Code{http: 505, grpc: grpcCodes.Unimplemented}
var VariantAlsoNegotiatesExperimental = Code{http: 506, grpc: grpcCodes.FailedPrecondition}
var LoopDetected = Code{http: 508, grpc: grpcCodes.Internal}
var NotExtended = Code{http: 510, grpc: grpcCodes.FailedPrecondition}
var NetworkAuthenticationRequired = Code{http: 511, grpc: grpcCodes.FailedPrecondition}
