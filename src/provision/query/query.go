package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var ec2Client *ec2.Client

func init() {

	// Create launch template
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return errorResponse(500, "Failed to load AWS config")
	}

	ec2Client = ec2.NewFromConfig(cfg)
}

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	// Log the incoming request
	log.Printf("Received request: %+v", request)

	// Parse query parameters
	queryParams := request.QueryStringParameters

	// Handle request body
	var body map[string]interface{}
	if request.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			log.Printf("Error decoding base64 body: %v", err)
			return errorResponse(http.StatusBadRequest, "Invalid base64 encoded body")
		}
		err = json.Unmarshal(decodedBody, &body)
	} else {
		err := json.Unmarshal([]byte(request.Body), &body)
		if err != nil {
			log.Printf("Error unmarshalling JSON body: %v", err)
			return errorResponse(http.StatusBadRequest, "Invalid JSON body")
		}
	}

	// Process the request based on the HTTP method
	switch request.HTTPMethod {
	case "GET":
		return handleGet(ctx, queryParams)
	case "POST":
		return handlePost(ctx, queryParams, body)
	case "DELETE":
		return handleDelete(ctx, queryParams, body)
	default:
		log.Printf("Method not allowed")
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

func main() {
	lambda.Start(handleRequest)
}
