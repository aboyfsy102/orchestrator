// I would like to write an AWS Lambda function that listens to ALB events and then writes to a S3 bucket.

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {

	return events.ALBTargetGroupResponse{
		StatusCode: 200,
	}, nil
}
