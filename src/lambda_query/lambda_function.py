import json
import logging
import boto3
import traceback

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def lambda_handler(event, context):
    # Log the event for debugging
    logger.info(f"Received event: {json.dumps(event)}")

    try:
        # Parse the ALB event
        # alb_event = event['Records'][0]['cf']['request']
        alb_event = event
        http_method = alb_event['httpMethod']
        path = alb_event['path']
        query_string_params = alb_event.get('queryStringParameters', {})
        headers = alb_event['headers']
        body = alb_event.get('body', '')

        # Process the request based on the HTTP method
        if http_method == 'GET':
            response_body = handle_get(path, query_string_params)
        elif http_method == 'POST':
            response_body = handle_post(path, json.loads(body) if body else {})
        elif http_method == 'DELETE':
            response_body = handle_delete(path, query_string_params)
        else:
            response_body = {"error": "Unsupported HTTP method"}
            return create_response(405, response_body)

        # Prepare the success response
        return create_response(200, response_body)

    except Exception as e:
        logger.error(f"Error processing request: {str(e)}\nStack trace:\n{traceback.format_exc()}")
        return create_response(500, {"error": "Internal Server Error"})

def handle_get(path, query_params):
    # Implement your GET logic here
    return {
        "message": "GET request processed",
        "path": path,
        "query_params": query_params
    }

def handle_post(path, body):
    try:
        # Create a boto3 client for EC2
        ec2_client = boto3.client('ec2')

        # Prepare the Spot Fleet request
        spot_fleet_request = {
            "SpotFleetRequestConfig": {
                "AllocationStrategy": body.get("AllocationStrategy", "lowestPrice"),
                "TargetCapacity": body["TargetCapacity"],
                "IamFleetRole": body["IamFleetRole"],
                "LaunchSpecifications": body["LaunchSpecifications"],
                "SpotPrice": body.get("SpotPrice", "0.03"),  # Default to $0.03 if not specified
                "TerminateInstancesWithExpiration": True,
                "Type": "request",
                "ReplaceUnhealthyInstances": False,
                "InstanceInterruptionBehavior": "terminate"
            }
        }

        # Create the Spot Fleet request
        response = ec2_client.request_spot_fleet(
            SpotFleetRequestConfig=spot_fleet_request["SpotFleetRequestConfig"]
        )

        return {
            "message": "Spot Fleet request created successfully",
            "SpotFleetRequestId": response["SpotFleetRequestId"]
        }

    except Exception as e:
        logger.error(f"Error creating Spot Fleet: {str(e)}")
        return {
            "message": "Error creating Spot Fleet",
            "error": str(e)
        }

def handle_delete(path, query_params):
    # Implement your DELETE logic here
    return {
        "message": "DELETE request processed",
        "path": path,
        "query_params": query_params
    }

def create_response(status_code, body):
    return {
        "statusCode": status_code,
        "statusDescription": f"{status_code} {get_status_description(status_code)}",
        "isBase64Encoded": False,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps(body)
    }

def get_status_description(status_code):
    descriptions = {
        200: "OK",
        405: "Method Not Allowed",
        500: "Internal Server Error"
    }
    return descriptions.get(status_code, "")