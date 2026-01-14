package errors

import (
	"sync"

	grpcCodes "google.golang.org/grpc/codes"
)

var stackTraceDepth = 50
var stackTraceDepthSetOnce sync.Once

func SetStackTraceDepth(depth int) {
	if depth < 1 {
		return
	}

	stackTraceDepthSetOnce.Do(func() {
		stackTraceDepth = depth
	})
}

var httpHeaderKey = "x-http-status"
var httpHeaderKeySetOnce sync.Once

func SetHttpHeaderKey(key string) {
	httpHeaderKeySetOnce.Do(func() {
		httpHeaderKey = key
	})
}

var defaultErrorCode = BadRequest
var defaultErrorCodeSetOnce sync.Once

func SetDefaultErrorCode(code Code) {
	defaultErrorCodeSetOnce.Do(func() {
		defaultErrorCode = code
	})
}

var defaultErrorCodeIfNotInMapper = BadRequest
var defaultCodeIfNotInMapperSetOnce sync.Once

func SetDefaultCodeIfNotInMapper(code Code) {
	defaultCodeIfNotInMapperSetOnce.Do(func() {
		defaultErrorCodeIfNotInMapper = code
	})
}

var codeMapperSetOnce sync.Once
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

func SetCodeMapper(mapper map[int]Code) {
	codeMapperSetOnce.Do(func() {
		codeMapper = mapper
	})
}
