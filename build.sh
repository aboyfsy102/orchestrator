#!/bin/bash

# Clean up
rm -rf dist
mkdir dist

# Build the lambda_eni_http
cd src/eni/http
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap .
zip ../../../dist/lambda_eni_http.zip bootstrap
rm bootstrap
cd ../../../

# Build the lambda_describe_http
cd src/describe/http
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap .
zip ../../../dist/lambda_describe_http.zip bootstrap
rm bootstrap
cd ../../../

# Run terraform
cd infra
# terraform init
terraform plan
terraform apply -auto-approve
cd ../