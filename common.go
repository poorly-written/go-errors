package errors

import (
	"encoding/json"
	"sync"

	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/anypb"
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

type extractDetails func(idx int, err any) (message *string, internalCode *string, reasons map[string][]Reason)

var detailsExtractor extractDetails = func(_ int, err any) (*string, *string, map[string][]Reason) {
	anyErr, ok := err.(*anypb.Any)
	if !ok {
		return nil, nil, nil
	}

	dErr := &DetailedErrorResponse{}
	if err := anyErr.UnmarshalTo(dErr); err != nil {
		return nil, nil, nil
	}

	var message, code *string
	if dErr.Message != "" {
		message = &dErr.Message
	}

	if dErr.Code != nil {
		code = dErr.Code
	}

	var reasons map[string][]Reason
	if dErr.Reasons != nil {
		r := make(map[string][]reason)
		if bytes, err := dErr.Reasons.MarshalJSON(); err == nil {
			if err := json.Unmarshal(bytes, &r); err == nil {
				for k, items := range r {
					if _, ok := reasons[k]; !ok {
						reasons[k] = make([]Reason, len(items))
					}

					for _, item := range items {
						reasons[k] = append(reasons[k], item)
					}
				}
			}
		}
	}

	return message, code, reasons
}

var reasonExtractorSetOnce sync.Once

func SetReasonExtractor(extractor extractDetails) {
	reasonExtractorSetOnce.Do(func() {
		detailsExtractor = extractor
	})
}
