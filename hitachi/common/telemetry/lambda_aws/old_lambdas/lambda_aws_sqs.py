import json
import boto3
import logging
import os
from datetime import datetime, timezone

logger = logging.getLogger()
logger.setLevel(logging.INFO)
# Initialize SQS
sqs = boto3.client("sqs")
SQS_QUEUE_URL = os.getenv("SQS_QUEUE_URL", "")

RESPONSE_MSG = {
    "statusCode": 200,
    "body": json.dumps({"message": "Request received, processing asynchronously."}),
}

RESPONSE_MSG_ERR = {
    "statusCode": 400,
    "body": json.dumps({"message": "Validation failed"}),
}

def lambda_handler(event, context):
    try:
        # Parse request body
        body = event
        site_id = body.get("site_id")

        if not site_id:
            return {"statusCode": 400, "body": json.dumps({"error": "Missing 'site_id'"})}

        current_time = datetime.now(timezone.utc).isoformat() + "Z"
        body["current_time"] = current_time

        if not validate_input(body):
            logger.error("body data validation failed")
            return RESPONSE_MSG_ERR

        # Send message to SQS
        response = sqs.send_message(
            QueueUrl=SQS_QUEUE_URL, MessageBody=json.dumps(body)
        )

        logger.info(f"Sent message to SQS: {response['MessageId']}")

        return RESPONSE_MSG

    except Exception as e:
        logger.error(f"Error: {e}")
        return {
            "statusCode": 500,
            "body": json.dumps({"error": "Failed to process request"}),
        }

def validate_input(data):
    schema = {
        "module_name": str,
        "operation_name": str,
        "operation_status": (int, lambda x: x in [0, 1]),
        "storage_model": str,
        "storage_serial": str,
        "connection_type": str,
        "storage_type": str,
        "process_time": ((int, float), lambda x: x > 0),  # Must be positive
        "site_id": (str, lambda x: len(x) < 121),
    }

    for key, expected in schema.items():
        if key not in data:
            logger.error("Required field missing in request data")
            return False
        value = data[key]
        if isinstance(expected, tuple):
            expected_type, validator = expected
            if not isinstance(value, expected_type) or not validator(value):
                logger.error("Validation failed for a field")
                return False
        else:
            if not isinstance(value, expected):
                logger.error("Validation failed for a field")
                return False

    return True
