package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/propagation"
)

var (
	_ sqsiface.SQSAPI = SQS{}
)

type SQS struct {
	*awssqs.SQS
	tracer   *zipkin.Tracer
	injector func(input *awssqs.SendMessageInput) propagation.Injector
}

func Wrap(s *awssqs.SQS, tracer *zipkin.Tracer) sqsiface.SQSAPI {
	return &SQS{
		SQS:      s,
		tracer:   tracer,
		injector: InjectSQS,
	}
}

func (s SQS) SendMessageWithContext(ctx aws.Context, input *awssqs.SendMessageInput, opts ...request.Option) (*awssqs.SendMessageOutput, error) {
	span := zipkin.SpanFromContext(ctx)
	if span == nil {
		return s.SQS.SendMessageWithContext(ctx, input, opts...)
	}

	err := s.injector(input)(span.Context())
	if err != nil {
		return s.SQS.SendMessageWithContext(ctx, input, opts...)
	}

	return s.SQS.SendMessageWithContext(ctx, input, opts...)
}
