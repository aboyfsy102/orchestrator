package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEC2Client is a mock of EC2ClientAPI
type MockEC2Client struct {
	mock.Mock
}

func (m *MockEC2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
}

func TestHandleRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        events.ALBTargetGroupRequest
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "GET request",
			request: events.ALBTargetGroupRequest{
				HTTPMethod: "GET",
			},
			expectedStatus: 200,
			expectedBody:   `[{"InstanceId":"i-1234567890abcdef0","InstanceType":"t2.micro","State":"running","PublicIpAddress":"203.0.113.1","PrivateIpAddress":"10.0.0.1"}]`,
		},
		{
			name: "Unsupported method",
			request: events.ALBTargetGroupRequest{
				HTTPMethod: "POST",
			},
			expectedStatus: 405,
			expectedBody:   `{"error": "Method not allowed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockEC2Client)
			if tt.request.HTTPMethod == "GET" {
				mockClient.On("DescribeInstances", mock.Anything, mock.Anything).Return(&ec2.DescribeInstancesOutput{
					Reservations: []types.Reservation{
						{
							Instances: []types.Instance{
								{
									InstanceId:       aws.String("i-1234567890abcdef0"),
									InstanceType:     types.InstanceTypeT2Micro,
									State:            &types.InstanceState{Name: types.InstanceStateNameRunning},
									PublicIpAddress:  aws.String("203.0.113.1"),
									PrivateIpAddress: aws.String("10.0.0.1"),
								},
							},
						},
					},
				}, nil)
			}

			// In the TestHandleRequest function, before calling handleRequest:
			ec2Client = mockClient

			// Then call handleRequest as before
			response, err := handleRequest(context.Background(), tt.request, mockClient)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.StatusCode)

			var responseBody, expectedBody interface{}
			json.Unmarshal([]byte(response.Body), &responseBody)
			json.Unmarshal([]byte(tt.expectedBody), &expectedBody)
			assert.Equal(t, expectedBody, responseBody)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestHandleGet(t *testing.T) {
	mockClient := new(MockEC2Client)
	mockClient.On("DescribeInstances", mock.Anything, mock.Anything).Return(&ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId:       aws.String("i-1234567890abcdef0"),
						InstanceType:     types.InstanceTypeT2Micro,
						State:            &types.InstanceState{Name: types.InstanceStateNameRunning},
						PublicIpAddress:  aws.String("203.0.113.1"),
						PrivateIpAddress: aws.String("10.0.0.1"),
					},
				},
			},
		},
	}, nil)

	response, err := handleGet(context.Background(), mockClient)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	var responseBody []map[string]string
	json.Unmarshal([]byte(response.Body), &responseBody)

	expectedBody := []map[string]string{
		{
			"InstanceId":       "i-1234567890abcdef0",
			"InstanceType":     "t2.micro",
			"State":            "running",
			"PublicIpAddress":  "203.0.113.1",
			"PrivateIpAddress": "10.0.0.1",
		},
	}

	assert.Equal(t, expectedBody, responseBody)

	mockClient.AssertExpectations(t)
}
