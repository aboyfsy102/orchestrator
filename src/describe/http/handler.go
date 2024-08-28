package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2ClientAPI defines the interface for EC2 client methods we're using
type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// Add this line to make the EC2 client interface accessible
var ec2Client EC2ClientAPI

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		// Use the ec2Client here instead of creating a new one
		return handleGet(ctx, ec2Client)
	default:
		return events.ALBTargetGroupResponse{
			StatusCode: 405,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func handleGet(ctx context.Context, client EC2ClientAPI) (events.ALBTargetGroupResponse, error) {
	result, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
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

	body, err := json.Marshal(instances)
	if err != nil {
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
	lambda.Start(handleRequest)
}
