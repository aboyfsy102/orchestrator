package llib

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func CreateError500Response(requestID string, err error) events.ALBTargetGroupResponse {
	body, _ := json.Marshal(map[string]string{
		"error":     err.Error(),
		"requestID": requestID,
	})

	return events.ALBTargetGroupResponse{
		StatusCode:        500,
		StatusDescription: "500 Internal Server Error",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Request-ID": requestID,
		},
		Body: string(body),
	}
}

func CreateSuccessResponse(requestID string, body string) events.ALBTargetGroupResponse {
	return events.ALBTargetGroupResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Request-ID": requestID,
		},
		Body: body,
	}
}
