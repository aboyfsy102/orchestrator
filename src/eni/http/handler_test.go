package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Update the mock to implement the new interface
type mockEC2Client struct {
	mock.Mock
}

// Implement the interface methods for the mock
func (m *mockEC2Client) DescribeNetworkInterfaces(ctx context.Context, params *ec2.DescribeNetworkInterfacesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeNetworkInterfacesOutput), args.Error(1)
}

func (m *mockEC2Client) CreateNetworkInterface(ctx context.Context, params *ec2.CreateNetworkInterfaceInput, optFns ...func(*ec2.Options)) (*ec2.CreateNetworkInterfaceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.CreateNetworkInterfaceOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteNetworkInterface(ctx context.Context, params *ec2.DeleteNetworkInterfaceInput, optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkInterfaceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteNetworkInterfaceOutput), args.Error(1)
}

// ... rest of the test file remains the same ...

func TestHandleGet(t *testing.T) {
	mockClient := new(mockEC2Client)
	mockClient.On("DescribeNetworkInterfaces", mock.Anything, mock.Anything).Return(&ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{
			{NetworkInterfaceId: stringPtr("eni-12345"), PrivateIpAddress: stringPtr("10.0.0.1")},
		},
	}, nil)

	resp, err := handleGet(context.Background(), mockClient)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body []map[string]string
	json.Unmarshal([]byte(resp.Body), &body)
	assert.Len(t, body, 1)
	assert.Equal(t, "eni-12345", body[0]["NetworkInterfaceId"])
	assert.Equal(t, "10.0.0.1", body[0]["PrivateIpAddress"])

	mockClient.AssertExpectations(t)
}

func TestHandlePost(t *testing.T) {
	mockClient := new(mockEC2Client)
	mockClient.On("CreateNetworkInterface", mock.Anything, mock.Anything).Return(&ec2.CreateNetworkInterfaceOutput{
		NetworkInterface: &types.NetworkInterface{NetworkInterfaceId: stringPtr("eni-67890")},
	}, nil)

	request := events.ALBTargetGroupRequest{
		HTTPMethod: "POST",
		Body:       `{"ip_address": "10.0.0.2"}`,
	}

	resp, err := handlePost(context.Background(), mockClient, request)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Body, "eni-67890")

	mockClient.AssertExpectations(t)
}

func TestHandleDelete(t *testing.T) {
	mockClient := new(mockEC2Client)
	mockClient.On("DeleteNetworkInterface", mock.Anything, mock.Anything).Return(&ec2.DeleteNetworkInterfaceOutput{}, nil)

	request := events.ALBTargetGroupRequest{
		HTTPMethod: "DELETE",
		Body:       `{"eni_id": "eni-12345"}`,
	}

	resp, err := handleDelete(context.Background(), mockClient, request)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Body, "deleted successfully")

	mockClient.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
