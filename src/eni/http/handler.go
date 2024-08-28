package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2ClientAPI defines the interface for EC2 client methods we're using
type EC2ClientAPI interface {
	DescribeNetworkInterfaces(ctx context.Context, params *ec2.DescribeNetworkInterfacesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error)
	CreateNetworkInterface(ctx context.Context, params *ec2.CreateNetworkInterfaceInput, optFns ...func(*ec2.Options)) (*ec2.CreateNetworkInterfaceOutput, error)
	DeleteNetworkInterface(ctx context.Context, params *ec2.DeleteNetworkInterfaceInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkInterfaceOutput, error)
}

func handleRequest(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.ALBTargetGroupResponse{
			StatusCode:        500,
			StatusDescription: "500 Internal Server Error",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            fmt.Sprintf(`{"error": "Error loading AWS config: %s"}`, err.Error()),
			IsBase64Encoded: false,
		}, err
	}

	client := ec2.NewFromConfig(cfg)

	switch request.HTTPMethod {
	case "GET":
		return handleGet(ctx, client)
	case "POST":
		return handlePost(ctx, client, request)
	case "DELETE":
		return handleDelete(ctx, client, request)
	default:
		return events.ALBTargetGroupResponse{
			StatusCode:        405,
			StatusDescription: "405 Method Not Allowed",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            `{"error": "Method not allowed"}`,
			IsBase64Encoded: false,
		}, nil
	}
}

func handleGet(ctx context.Context, client EC2ClientAPI) (events.ALBTargetGroupResponse, error) {
	result, err := client.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{})
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

	interfaces := make([]map[string]string, len(result.NetworkInterfaces))
	for i, eni := range result.NetworkInterfaces {
		interfaces[i] = map[string]string{
			"NetworkInterfaceId": *eni.NetworkInterfaceId,
			"PrivateIpAddress":   *eni.PrivateIpAddress,
		}
	}

	body, err := json.Marshal(interfaces)
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

func handlePost(ctx context.Context, client EC2ClientAPI, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	var reqBody struct {
		IPAddress string `json:"ip_address"`
	}
	err := json.Unmarshal([]byte(request.Body), &reqBody)
	if err != nil {
		return events.ALBTargetGroupResponse{
			StatusCode:        400,
			StatusDescription: "400 Bad Request",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            fmt.Sprintf(`{"error": "Invalid request body: %s"}`, err.Error()),
			IsBase64Encoded: false,
		}, err
	}

	result, err := client.CreateNetworkInterface(ctx, &ec2.CreateNetworkInterfaceInput{
		PrivateIpAddress: &reqBody.IPAddress,
	})
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

	return events.ALBTargetGroupResponse{
		StatusCode:        200,
		StatusDescription: "200 OK",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            fmt.Sprintf(`{"message": "Network interface created", "eni_id": "%s"}`, *result.NetworkInterface.NetworkInterfaceId),
		IsBase64Encoded: false,
	}, nil
}

func handleDelete(ctx context.Context, client EC2ClientAPI, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	var reqBody struct {
		ENIID string `json:"eni_id"`
	}
	err := json.Unmarshal([]byte(request.Body), &reqBody)
	if err != nil {
		return events.ALBTargetGroupResponse{
			StatusCode:        400,
			StatusDescription: "400 Bad Request",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:            fmt.Sprintf(`{"error": "Invalid request body: %s"}`, err.Error()),
			IsBase64Encoded: false,
		}, err
	}

	_, err = client.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: &reqBody.ENIID,
	})
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

	return events.ALBTargetGroupResponse{
		StatusCode:        200,
		StatusDescription: "200 OK",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            fmt.Sprintf(`{"message": "Network interface %s deleted successfully"}`, reqBody.ENIID),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
