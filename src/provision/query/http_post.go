package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func handlePost(ctx context.Context, cwl *cloudwatchlogs.Client, queryParams map[string]string, body map[string]interface{}) (events.ALBTargetGroupResponse, error) {
	// Process POST request
	response := fmt.Sprintf("Received POST request with query parameters: %v and body: %v", queryParams, body)
	logToCloudWatch(ctx, cwl, response)
	return successResponse(response)
}
