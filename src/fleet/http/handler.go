package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2ClientAPI defines the interface for EC2 client methods we're using
type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest, client EC2ClientAPI) (events.ALBTargetGroupResponse, error) {
	// Get Lambda context for request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := lc.AwsRequestID

	log.Printf("Received request: Method=%s, Path=%s", request.HTTPMethod, request.Path)

	switch request.HTTPMethod {
	case "GET":
		log.Println("Handling GET request")
		return handleGet(ctx, client, requestID)
	case "POST":
		log.Println("Handling POST request")
		return handlePost(ctx, client, requestID)
	case "DELETE":
		log.Println("Handling DELETE request")
		return handleDelete(ctx, client, requestID)
	default:
		log.Printf("Unsupported method: %s", request.HTTPMethod)
		return events.ALBTargetGroupResponse{
			StatusCode: 405,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func handleGet(ctx context.Context, client EC2ClientAPI, requestID string) (events.ALBTargetGroupResponse, error) {
	// TODO: Implement POST request handling
	return events.ALBTargetGroupResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented"}`,
	}, nil
}

func handlePost(ctx context.Context, client EC2ClientAPI, requestID string) (events.ALBTargetGroupResponse, error) {
	// TODO: Implement POST request handling
	return events.ALBTargetGroupResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented"}`,
	}, nil
}

func handleDelete(ctx context.Context, client EC2ClientAPI, requestID string) (events.ALBTargetGroupResponse, error) {
	// TODO: Implement POST request handling
	return events.ALBTargetGroupResponse{
		StatusCode: 501,
		Body:       `{"error": "Not implemented"}`,
	}, nil
}

func main() {
	log.Println("Lambda function starting")

	log.Println("Initializing AWS SDK configuration")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	log.Println("Creating EC2 client")
	ec2Client := ec2.NewFromConfig(cfg)

	log.Println("Starting Lambda handler")
	lambda.Start(func(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
		return handleRequest(ctx, request, ec2Client)
	})
}
