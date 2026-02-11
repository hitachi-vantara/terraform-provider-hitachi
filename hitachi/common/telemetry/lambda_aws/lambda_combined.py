import json
import boto3
import logging
import os
import time
import urllib3
from datetime import datetime, timezone
from botocore.auth import SigV4Auth
from botocore.awsrequest import AWSRequest

logger = logging.getLogger()
logger.setLevel(logging.INFO)

# ---- Environment Variables ----
OPENSEARCH_ENDPOINT = os.getenv("OPENSEARCH_ENDPOINT", "")
REGION = os.getenv("AWS_REGION", "us-west-2")   # FIX: ensure REGION is never None
INDEX = "terraform_telemetry_v2"

# ---- Static headers ----
HEADERS = {"Content-Type": "application/json", "user-agent": "terraform"}

http = urllib3.PoolManager()

# --------------------------------------------------------------------------
#  SIGV4 SIGNED REQUEST
# --------------------------------------------------------------------------

def signed_request(method, url, body=None, headers=None):
    """Send a SigV4-signed HTTP request to OpenSearch."""
    session = boto3.Session()

    creds = session.get_credentials()
    region = session.region_name or REGION  # FIX: fallback to REGION

    aws_request = AWSRequest(method=method, url=url, data=body, headers=headers or {})
    SigV4Auth(creds, "es", region).add_auth(aws_request)

    return http.request(
        method,
        url,
        body=body,
        headers=dict(aws_request.headers),
    )

# --------------------------------------------------------------------------
#  VALIDATION
# --------------------------------------------------------------------------

def validate_input(data):
    schema = {
        "module_name": str,
        "operation_name": str,
        "operation_status": (int, lambda x: x in [0, 1]),
        "storage_model": str,
        "storage_serial": str,
        "connection_type": str,
        "storage_type": str,
        "process_time": ((int, float), lambda x: x > 0),
        "site_id": (str, lambda x: len(x) < 121),
    }

    for key, expected in schema.items():
        if key not in data:
            logger.error(f"Missing key: {key}")
            return False

        value = data[key]

        if isinstance(expected, tuple):
            expected_type, validator = expected
            if not isinstance(value, expected_type) or not validator(value):
                logger.error(f"Invalid value for {key}: {value}")
                return False
        else:
            if not isinstance(value, expected):
                logger.error(f"Incorrect type for {key}: {value}")
                return False

    return True

# --------------------------------------------------------------------------
#  OPENSEARCH UPSERT (unchanged except region fix)
# --------------------------------------------------------------------------

def update_or_create_document(site_id, body, max_retries=3, retry_delay=1):
    current_time = body["current_time"]
    success = body["operation_status"]
    failure = 1 - success
    module_name = body["module_name"]
    operation_name = body["operation_name"]
    process_time = body.get("process_time", 0)
    storage_model = body["storage_model"]
    storage_serial = body["storage_serial"]
    storage_type = body["storage_type"]
    connection_type = body["connection_type"]

    task_key_name = f"{module_name}.{operation_name}"

    script_source = """
    def storage_type = params.storage_type;
    if (ctx._source[storage_type] == null) {
        ctx._source[storage_type] = [];
    }
    def storage = ctx._source[storage_type].find(s -> s.model == params.storage_model && s.serial == params.storage_serial);
    if (storage == null) {
        storage = ['model': params.storage_model, 'serial': params.storage_serial, params.connection_type: []];
        ctx._source[storage_type].add(storage);
    }
    if (storage[params.connection_type] == null) {
        storage[params.connection_type] = [];
    }
    def existing = storage[params.connection_type].find(t -> t.name == params.task_key_name);
    if (existing == null) {
        def newTask = ['name': params.task_key_name, 'metrics': ['success': params.success, 'failure': params.failure, 'averageTimeInSec': params.process_time]];
        storage[params.connection_type].add(newTask);
    } else {
        def metrics = existing.metrics;
        metrics.success += params.success;
        metrics.failure += params.failure;
        def total = metrics.success + metrics.failure;
        metrics.averageTimeInSec = Math.round(((metrics.averageTimeInSec * (total - 1)) + params.process_time) / total * 100.0) / 100.0;
    }
    ctx._source.lastUpdate = params.current_time;
    """

    script = {
        "source": script_source,
        "lang": "painless",
        "params": {
            "storage_model": storage_model,
            "storage_serial": storage_serial,
            "task_key_name": task_key_name,
            "process_time": process_time,
            "success": success,
            "failure": failure,
            "connection_type": connection_type,
            "current_time": current_time,
            "storage_type": storage_type,
        },
    }

    upsert_doc = {
        "siteId": site_id,
        "createDate": current_time,
        "lastUpdate": current_time,
        storage_type: [
            {
                "model": storage_model,
                "serial": storage_serial,
                connection_type: [
                    {
                        "name": task_key_name,
                        "metrics": {
                            "success": success,
                            "failure": failure,
                            "averageTimeInSec": process_time,
                        }
                    }
                ],
            }
        ],
    }

    update_url = f"{OPENSEARCH_ENDPOINT}/{INDEX}/_update/{site_id}?retry_on_conflict=5"
    payload = json.dumps({"script": script, "upsert": upsert_doc})

    for attempt in range(max_retries):
        try:
            response = signed_request("POST", update_url, body=payload, headers=HEADERS)

            logger.info("OpenSearch payload: %s", payload)

            if response.status in (200, 201):
                logger.info(f"Successfully processed site_id {site_id}")
                return json.loads(response.data.decode())

            if response.status == 409:
                logger.warning(f"409 conflict attempt {attempt+1}, retrying...")
                time.sleep(retry_delay)
                retry_delay *= 2
                continue

            raise Exception(f"OpenSearch error {response.status}: {response.data}")

        except Exception as e:
            logger.error(f"Attempt {attempt+1} failed: {e}")
            time.sleep(retry_delay)
            retry_delay *= 2

    raise Exception("Failed after max retries")

# --------------------------------------------------------------------------
#  COMBINED HANDLER
# --------------------------------------------------------------------------

def lambda_handler(event, context):
    try:
        # If coming from API Gateway, parse the JSON body string
        if isinstance(event, dict) and "body" in event:
            body = json.loads(event["body"] or "{}")
        else:
            body = event if isinstance(event, dict) else json.loads(event)

        site_id = body.get("site_id")
        if not site_id:
            return {"statusCode": 400, "body": json.dumps({"error": "Missing 'site_id'"})}

        body["current_time"] = datetime.now(timezone.utc).isoformat() + "Z"

        if not validate_input(body):
            logger.error(f"Validation failed for body: {body}")
            return {"statusCode": 400, "body": json.dumps({"error": "Validation failed"})}

        if not OPENSEARCH_ENDPOINT:
            logger.warning("OPENSEARCH_ENDPOINT not set.")
            return {
                "statusCode": 500,
                "body": json.dumps({"message": "Failed: OPENSEARCH_ENDPOINT not set."})
            }

        # Write to OpenSearch
        update_or_create_document(site_id, body)

        return {
            "statusCode": 200,
            "body": json.dumps({"message": "Processed successfully"})
        }

    except Exception as e:
        logger.error(f"Error: {str(e)}", exc_info=True)
        return {
            "statusCode": 500,
            "body": json.dumps({"error": "Failed to process request"})
        }
