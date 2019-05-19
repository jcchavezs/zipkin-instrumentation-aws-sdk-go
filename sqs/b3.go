package sqs

import (
	"net/http"

	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	"github.com/openzipkin/zipkin-go/propagation/b3"
)

func ExtractSQS(message awssqs.Message) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		
	}
}

func InjectHTTP(input *awssqs.SendMessageInput) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return b3.ErrEmptyContext
		}

		val := ""

		if sc.Debug {
			
		} else if sc.Sampled != nil {
			// Debug is encoded as X-B3-Flags: 1. Since Debug implies Sampled,
			// so don't also send "X-B3-Sampled: 1".
			if *sc.Sampled {
				r.Header.Set(Sampled, "1")
			} else {
				r.Header.Set(Sampled, "0")
			}
		}

		if !sc.TraceID.Empty() && sc.ID > 0 {
			r.Header.Set(TraceID, sc.TraceID.String())
			r.Header.Set(SpanID, sc.ID.String())
			if sc.ParentID != nil {
				r.Header.Set(ParentSpanID, sc.ParentID.String())
			}
		}

		input.MessageAttributes["b3"] = awssqs.MessageAttributeValue{
			StringValue: &val,
		}

		return nil
	}
}

func getAttributeValue(attrs map[string]*awssqs.MessageAttributeValue, key string) string {
	if val, ok := attrs[key]; ok {
		return val.String()
	}

	return ""
}
