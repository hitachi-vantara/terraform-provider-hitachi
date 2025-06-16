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
    "status": 200,
    "body": json.dumps({"message": "Request received, processing asynchronously."}),
}

def lambda_handler(event, context):
    try:
        # Parse request body
        body = event
        site_id = body.get("site")

        if not site_id:
            return {"status": 400, "body": json.dumps({"error": "Missing 'site_id'"})}

        current_time = datetime.now(timezone.utc).isoformat() + "Z"
        body["current_time"] = current_time

        if not validate_input(body):
            logger.error(f"body data validation failed: {body}")
            return RESPONSE_MSG

        # Send message to SQS
        response = sqs.send_message(
            QueueUrl=SQS_QUEUE_URL, MessageBody=json.dumps(body)
        )

        logger.info(f"Sent message to SQS: {response['MessageId']}")
        logger.info(f"Payload: {json.dumps(body)}")

        return RESPONSE_MSG

    except Exception as e:
        logger.error(f"Error: {e}", exc_info=True)
        return {
            "status": 500,
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
    }

    for key, expected in schema.items():
        if key not in data:
            return False
        value = data[key]
        if isinstance(expected, tuple):
            expected_type, validator = expected
            if not isinstance(value, expected_type) or not validator(value):
                return False
        else:
            if not isinstance(value, expected):
                return False

    return True
