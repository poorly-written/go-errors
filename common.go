package errors

import (
	"context"
	"encoding/json"
	"sync"

	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
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

var defaultErrorCode Code = BadRequest
var defaultErrorCodeSetOnce sync.Once

func SetDefaultErrorCode(code Code) {
	defaultErrorCodeSetOnce.Do(func() {
		defaultErrorCode = code
	})
}

var defaultErrorCodeIfNotInMapper Code = BadRequest
var defaultCodeIfNotInMapperSetOnce sync.Once

func SetDefaultCodeIfNotInMapper(code Code) {
	defaultCodeIfNotInMapperSetOnce.Do(func() {
		defaultErrorCodeIfNotInMapper = code
	})
}

var codeMapperSetOnce sync.Once
var useCodeMapperLock sync.Mutex
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

func AddToCodeMapper(when int, code Code) {
	useCodeMapperLock.Lock()
	defer useCodeMapperLock.Unlock()
	codeMapper[when] = code
}

type ErrorDetails struct {
	Message         *string
	InternalCode    *string
	Reasons         map[string][]Reason
	IncludeMetadata bool
	Metadata        map[string]interface{}
}

type errorUnmarshalerFunc func(idx int, err any) (*ErrorDetails, error)

var errorUnmarshaler errorUnmarshalerFunc = func(_ int, err any) (*ErrorDetails, error) {
	anyErr, ok := err.(*anypb.Any)
	if !ok {
		return nil, nil
	}

	dErr := &DetailedErrorResponse{}
	if err := anyErr.UnmarshalTo(dErr); err != nil {
		return nil, err
	}

	var details = &ErrorDetails{}

	if dErr.Message != "" {
		details.Message = &dErr.Message
	}

	if dErr.Code != nil {
		details.InternalCode = dErr.Code
	}

	if dErr.Metadata != nil {
		bytes, err := dErr.Metadata.MarshalJSON()
		if err != nil {
			return nil, err
		}

		var metadata map[string]interface{}
		if err := json.Unmarshal(bytes, &metadata); err != nil {
			return nil, err
		}

		details.Metadata = metadata
	}

	if dErr.Reasons != nil {
		bytes, err := dErr.Reasons.MarshalJSON()
		if err != nil {
			return nil, err
		}

		reasons := make(map[string][]reason)
		if err := json.Unmarshal(bytes, &reasons); err != nil {
			return nil, err
		}

		list := make(map[string][]Reason, len(reasons))
		for k, v := range reasons {
			list[k] = make([]Reason, len(v))
			for i, item := range v {
				list[k][i] = item
			}
		}

		details.Reasons = list
	}

	return details, nil
}

var errorUnmarshalerSetOnce sync.Once

func SetErrorUnmarshaler(unmarshaler errorUnmarshalerFunc) {
	errorUnmarshalerSetOnce.Do(func() {
		errorUnmarshaler = unmarshaler
	})
}

type errorMarshalerFunc func(*ErrorDetails) (*anypb.Any, error)

var errorMarshaler errorMarshalerFunc = func(details *ErrorDetails) (*anypb.Any, error) {
	// [Concept] https://stackoverflow.com/a/75720585/2190689
	de := &DetailedErrorResponse{
		Error: true,
	}

	if details.Message != nil {
		de.Message = *details.Message
	}

	if details.InternalCode != nil {
		de.Code = details.InternalCode
	}

	if reasons := details.Reasons; len(reasons) > 0 {
		reasonsMap := make(map[string]interface{})
		for key, items := range reasons {
			list := make([]interface{}, len(items))
			for i, each := range items {
				list[i] = each.ToHashMap()
			}

			reasonsMap[key] = list
		}

		reasonPb, err := structpb.NewStruct(reasonsMap)
		if err != nil {
			return nil, err
		}

		de.Reasons = reasonPb
	}

	if details.IncludeMetadata && len(details.Metadata) > 0 {
		mdPb, err := structpb.NewStruct(details.Metadata)
		if err != nil {
			return nil, err
		}

		de.Metadata = mdPb
	}

	return anypb.New(de)
}

var errorMarshalerSetOnce sync.Once

func SetErrorMarshaler(marshaler errorMarshalerFunc) {
	errorMarshalerSetOnce.Do(func() {
		errorMarshaler = marshaler
	})
}

type contextualMetadataExtractorFunc func(ctx context.Context) map[string]interface{}

var contextualMetadataExtractorSetOnce sync.Once
var contextualMetadataExtractor contextualMetadataExtractorFunc = func(ctx context.Context) map[string]interface{} {
	return nil
}

func SetContextualMetadataExtractor(extractor contextualMetadataExtractorFunc) {
	contextualMetadataExtractorSetOnce.Do(func() {
		contextualMetadataExtractor = extractor
	})
}
