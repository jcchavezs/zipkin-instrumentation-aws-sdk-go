package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	zipkin "github.com/openzipkin/zipkin-go"
)

var (
	_ sqsiface.SQSAPI = SQS{}
)

type SQS struct {
	parent *awssqs.SQS
	tracer *zipkin.Tracer
}

func (s *SQS) SendMessageWithContext(ctx aws.Context, input *awssqs.SendMessageInput) (*awssqs.SendMessageOutput, error) {
	req, out := s.parent.SendMessageRequest(input)
	return out, req.Send()
}
