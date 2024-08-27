package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func handlePost(ctx context.Context, queryParams map[string]string, body map[string]interface{}) (events.ALBTargetGroupResponse, error) {
	log.Printf("Handling POST request with query parameters: %v and body: %v", queryParams, body)

	// Create launch template
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return errorResponse(500, "Failed to load AWS config")
	}

	ec2Client := ec2.NewFromConfig(cfg)

	input := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: aws.String("MyLaunchTemplate"),
		VersionDescription: aws.String("Initial version"),
		LaunchTemplateData: &ec2.LaunchTemplateData{
			InstanceType: aws.String("t2.micro"),
			ImageId:      aws.String("ami-12345678"), // Replace with your desired AMI ID
			// Add more configuration options as needed
		},
	}

	result, err := ec2Client.CreateLaunchTemplate(ctx, input)
	if err != nil {
		return errorResponse(500, "Failed to create launch template: "+err.Error())
	}

	response := fmt.Sprintf("Created launch template: %s", *result.LaunchTemplate.LaunchTemplateId)

	return successResponse(response)
}
