import json
import time
import boto3
import logging
from botocore.exceptions import ClientError

logger = logging.getLogger()
logger.setLevel(logging.INFO)

# Initialize AWS clients
sqs_client = boto3.client('sqs')
ec2_client = boto3.client('ec2')


def lambda_handler(event, context):

    # Process each record from SQS
    for record in event['Records']:
        success = process_spot_fleet_request(record, sqs_client, ec2_client)
        if not success:
            return_message_to_queue(record)

    return {
        'statusCode': 200,
        'body': json.dumps('Processing complete')
    }

def process_spot_fleet_request(record):
    # Get the message body
    message_body = json.loads(record['body'])
    logger.info(f"Received message: {json.dumps(message_body)}")

    try:
        launch_template_id = create_launch_template(message_body, ec2_client)

        # Prepare the Fleet request
        fleet_request = {
            "Type": "instant",
            "TargetCapacitySpecification": {
                "TotalTargetCapacity": message_body["TargetCapacity"],
                "OnDemandTargetCapacity": 0,
                "DefaultTargetCapacityType": "spot"
            },
            "SpotOptions": {
                "AllocationStrategy": "price-capacity-optimized",
                "InstanceInterruptionBehavior": "terminate"
            },
            "LaunchTemplateConfigs": [
                {
                    "LaunchTemplateSpecification": {
                        "LaunchTemplateId": launch_template_id
                    },
                    "Overrides": [
                        {
                            "InstanceType": instance_type,
                            "SubnetId": subnet_id
                        }
                        for instance_type in spec.get("InstanceTypes", [])
                        for subnet_id in spec.get("SubnetIds", [])
                    ]
                }
                for spec in message_body["LaunchSpecifications"]
            ]
        }

        # Create the Fleet
        response = ec2_client.create_fleet(**fleet_request)

        logger.info(f"Fleet created successfully. Fleet ID: {response['FleetId']}")

        # Delete the message from the queue
        sqs_client.delete_message(
            QueueUrl=record['eventSourceARN'].split(':')[5],
            ReceiptHandle=record['receiptHandle']
        )

        return True  # Processing successful

    except ClientError as e:
        logger.error(f"Error creating Fleet: {e}")
    except KeyError as e:
        logger.error(f"Missing required key in message: {e}")
    except Exception as e:
        logger.error(f"Unexpected error: {e}")

    return False  # Processing failed

def return_message_to_queue(record):
    try:
        # Get the queue URL
        queue_url = record['eventSourceARN'].split(':')[5]

        # Send the message back to the queue
        sqs_client.send_message(
            QueueUrl=queue_url,
            MessageBody=record['body'],
            DelaySeconds=300  # 5 minutes delay before the message is available again
        )

        logger.info(f"Message returned to queue for retry: {record['messageId']}")

        # Delete the original message to prevent duplicate processing
        sqs_client.delete_message(
            QueueUrl=queue_url,
            ReceiptHandle=record['receiptHandle']
        )

    except Exception as e:
        logger.error(f"Error returning message to queue: {e}")


def create_launch_template(message_body):
    # Get the latest Amazon Linux 2023 AMI ID
    response = ec2_client.describe_images(
        Owners=['amazon'],
        Filters=[
            {'Name': 'name', 'Values': ['al2023-ami-*-x86_64']},
            {'Name': 'state', 'Values': ['available']}
        ]
    )
    ami_id = sorted(response['Images'], key=lambda x: x['CreationDate'], reverse=True)[0]['ImageId']

    # Prepare the launch template data
    launch_template_data = {
        'ImageId': ami_id,
        'InstanceType': message_body.get('InstanceType', 't3.micro'),
        'KeyName': message_body.get('KeyName'),
        'SecurityGroupIds': get_security_group_id('j5v3-'),
        'UserData': message_body.get('UserData'),
        'IamInstanceProfile': {
            'Name': 'ec2-instance-role',
        },
        'BlockDeviceMappings': [
            {
                'DeviceName': '/dev/xvda',
                'Ebs': {
                    'VolumeSize': message_body.get('VolumeSize', 8),
                    'VolumeType': 'gp3',
                    'DeleteOnTermination': True
                }
            }
        ],
        'TagSpecifications': [
            {
                'ResourceType': 'instance',
                'Tags': [
                    {'Key': 'Name', 'Value': message_body.get('InstanceName', 'SpotFleetInstance')},
                    {'Key': 'CreatedBy', 'Value': 'SpotFleetLambda'}
                ]
            }
        ]
    }

    # Create the launch template
    response = ec2_client.create_launch_template(
        LaunchTemplateName=f"SpotFleetTemplate-{message_body.get('InstanceName', 'Default')}-{int(time.time())}",
        VersionDescription='Initial version',
        LaunchTemplateData=launch_template_data
    )

    return response['LaunchTemplate']['LaunchTemplateId']


def get_security_group_id(security_group_name):
    response = ec2_client.describe_security_groups(
        Filters=[
            {'Name': 'group-name', 'Values': [security_group_name]}
        ]
    )
    return response['SecurityGroups'][0]['GroupId']