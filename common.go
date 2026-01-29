package errors

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/poorly-written/grpc-http-response/codes"
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

var defaultErrorCode codes.Code = codes.BadRequest
var defaultErrorCodeSetOnce sync.Once

func SetDefaultErrorCode(code codes.Code) {
	defaultErrorCodeSetOnce.Do(func() {
		defaultErrorCode = code
	})
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
