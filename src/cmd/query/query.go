package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// Initialize a global logger
var log = lambda.NewLogger()

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest, cwl *cloudwatchlogs.Client) (events.ALBTargetGroupResponse, error) {
	// Log the incoming request
	logToCloudWatch(ctx, cwl, fmt.Sprintf("Received request: %+v", request))

	// Parse query parameters
	queryParams := request.QueryStringParameters

	// Handle request body
	var body map[string]interface{}
	if request.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			logToCloudWatch(ctx, cwl, fmt.Sprintf("Error decoding base64 body: %v", err))
			return errorResponse(http.StatusBadRequest, "Invalid base64 encoded body")
		}
		err = json.Unmarshal(decodedBody, &body)
	} else {
		err := json.Unmarshal([]byte(request.Body), &body)
		if err != nil {
			logToCloudWatch(ctx, cwl, fmt.Sprintf("Error unmarshalling JSON body: %v", err))
			return errorResponse(http.StatusBadRequest, "Invalid JSON body")
		}
	}

	// Process the request based on the HTTP method
	switch request.HTTPMethod {
	case "GET":
		return handleGet(ctx, cwl, queryParams)
	case "POST":
		return handlePost(ctx, cwl, queryParams, body)
	case "DELETE":
		return handleDelete(ctx, cwl, queryParams, body)
	default:
		logToCloudWatch(ctx, cwl, "Method not allowed")
		return errorResponse(http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func successResponse(body string) (events.ALBTargetGroupResponse, error) {
	return events.ALBTargetGroupResponse{
		StatusCode:        200,
		StatusDescription: "200 OK",
		Headers:           map[string]string{"Content-Type": "application/json"},
		Body:              body,
		IsBase64Encoded:   false,
	}, nil
}

func errorResponse(statusCode int, message string) (events.ALBTargetGroupResponse, error) {
	return events.ALBTargetGroupResponse{
		StatusCode:        statusCode,
		StatusDescription: fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Headers:           map[string]string{"Content-Type": "application/json"},
		Body:              fmt.Sprintf(`{"error": "%s"}`, message),
		IsBase64Encoded:   false,
	}, nil
}

func logToCloudWatch(ctx context.Context, cwl *cloudwatchlogs.Client, message string) {
	lc, _ := lambdacontext.FromContext(ctx)
	logGroupName := os.Getenv("LOG_GROUP_NAME")
	if logGroupName == "" {
		log.Printf("LOG_GROUP_NAME environment variable not set")
		return
	}

	_, err := cwl.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(lc.AwsRequestID),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
			},
		},
	})
	if err != nil {
		log.Printf("Failed to log to CloudWatch: %v", err)
	}
}

func main() {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Failed to load AWS configuration: %v", err)
		return
	}

	// Create CloudWatch Logs client
	cwl := cloudwatchlogs.NewFromConfig(cfg)

	lambda.Start(func(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
		return handleRequest(ctx, request, cwl)
	})
}
