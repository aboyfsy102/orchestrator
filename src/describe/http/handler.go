package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2ClientAPI defines the interface for EC2 client methods we're using
type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// Add this line to make the EC2 client interface accessible
var ec2Client EC2ClientAPI

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest, client EC2ClientAPI) (events.ALBTargetGroupResponse, error) {
	log.Printf("Received request: Method=%s, Path=%s", request.HTTPMethod, request.Path)

	switch request.HTTPMethod {
	case "GET":
		log.Println("Handling GET request")
		return handleGet(ctx, client)
	default:
		log.Printf("Unsupported method: %s", request.HTTPMethod)
		return events.ALBTargetGroupResponse{
			StatusCode: 405,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func handleGet(ctx context.Context, client EC2ClientAPI) (events.ALBTargetGroupResponse, error) {
	log.Println("Starting DescribeInstances API call")
	result, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Printf("Error in DescribeInstances: %v", err)
		return events.ALBTargetGroupResponse{
			StatusCode:        500,
			StatusDescription: "500 Internal Server Error",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			IsBase64Encoded: false,
		}, err
	}
	log.Println("DescribeInstances API call completed successfully")

	log.Println("Processing instance data")
	instances := []map[string]string{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceInfo := map[string]string{
				"InstanceId":   *instance.InstanceId,
				"InstanceType": string(instance.InstanceType),
				"State":        string(instance.State.Name),
			}
			if instance.PublicIpAddress != nil {
				instanceInfo["PublicIpAddress"] = *instance.PublicIpAddress
			}
			if instance.PrivateIpAddress != nil {
				instanceInfo["PrivateIpAddress"] = *instance.PrivateIpAddress
			}
			instances = append(instances, instanceInfo)
		}
	}
	log.Printf("Processed %d instances", len(instances))

	log.Println("Marshaling response body")
	body, err := json.Marshal(instances)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return events.ALBTargetGroupResponse{
			StatusCode:        500,
			StatusDescription: "500 Internal Server Error",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            fmt.Sprintf(`{"error": "Failed to marshal response: %s"}`, err.Error()),
			IsBase64Encoded: false,
		}, err
	}

	log.Println("Returning successful response")
	return events.ALBTargetGroupResponse{
		StatusCode:        200,
		StatusDescription: "200 OK",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(body),
		IsBase64Encoded: false,
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
	ec2Client = ec2.NewFromConfig(cfg)

	log.Println("Starting Lambda handler")
	lambda.Start(func(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
		return handleRequest(ctx, request, ec2Client)
	})
}
