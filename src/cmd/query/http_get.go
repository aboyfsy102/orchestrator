package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func handleGet(ctx context.Context, cwl *cloudwatchlogs.Client, queryParams map[string]string) (events.ALBTargetGroupResponse, error) {
	// Process GET request
	response := fmt.Sprintf("Received GET request with query parameters: %v", queryParams)
	logToCloudWatch(ctx, cwl, response)
	return successResponse(response)
}
