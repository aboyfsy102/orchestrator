package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func handlePost(ctx context.Context, queryParams map[string]string, body map[string]interface{}) (events.ALBTargetGroupResponse, error) {
	panic("unimplemented")
}
