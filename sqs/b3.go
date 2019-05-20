package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	"github.com/openzipkin/zipkin-go/propagation/b3"
)

const b3Key = "b3"

// ExtractSQS extracts returns an extractor for a given sqs.Message
func ExtractSQS(message awssqs.Message) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		return b3.ParseSingleHeader(getB3AttributeValue(message.MessageAttributes))
	}
}

// InjectSQS returns an injector for a given sqs.SendMessageInput
func InjectSQS(input *awssqs.SendMessageInput) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return b3.ErrEmptyContext
		}

		input.MessageAttributes[b3Key] = &awssqs.MessageAttributeValue{
			StringValue: aws.String(b3.BuildSingleHeader(sc)),
			DataType:    aws.String("String"),
		}

		return nil
	}
}

func getB3AttributeValue(attrs map[string]*awssqs.MessageAttributeValue) string {
	if val, ok := attrs[b3Key]; ok {
		return val.String()
	}

	return ""
}
