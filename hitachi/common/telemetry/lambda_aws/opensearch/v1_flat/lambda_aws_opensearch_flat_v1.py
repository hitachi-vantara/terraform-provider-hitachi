import json
import boto3
import logging
import requests
from requests_aws4auth import AWS4Auth
import os
import time

logger = logging.getLogger()
logger.setLevel(logging.INFO)

OPENSEARCH_ENDPOINT = os.getenv("OPENSEARCH_ENDPOINT", "")
REGION = "us-west-2"
INDEX = "terraform_telemetry"
DOCUMENT_ID = "terraform_doc_id1"
SQS_QUEUE_URL = os.getenv("SQS_QUEUE_URL", "")

HEADERS = {"Content-Type": "application/json", "user-agent": "terraform"}

def get_awsauth():
    session = boto3.Session(region_name=REGION)
    credentials = session.get_credentials()
    if credentials is None:
        raise RuntimeError("AWS credentials not found")
    return AWS4Auth(
        credentials.access_key,
        credentials.secret_key,
        REGION,
        "es",
        session_token=credentials.token,
    )

def get_sqs_client():
    return boto3.client("sqs", region_name=REGION)

def update_or_create_document(site_id, body, max_retries=3, retry_delay=1):
    current_time = body.get("current_time")
    success = body.get("operation_status")
    failure = 1 - success
    module_name = body.get("module_name")
    operation_name = body.get("operation_name")
    process_time = body.get("process_time", 0)
    storage_model = body.get("storage_model")
    storage_serial = body.get("storage_serial")
    storage_type = body.get("storage_type")
    connection_type = body.get("connection_type")

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
    def task = storage[params.connection_type].find(t -> t.containsKey(params.task_key_name));
    if (task == null) {
        task = [params.task_key_name: ['success': params.success, 'failure': params.failure, 'averageTimeInSec': params.process_time]];
        storage[params.connection_type].add(task);
    } else {
        def existing = task[params.task_key_name];
        existing.success += params.success;
        existing.failure += params.failure;
        def total = existing.success + existing.failure;
        existing.averageTimeInSec = Math.round(((existing.averageTimeInSec * (total - 1)) + params.process_time) / total * 100.0) / 100.0;
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
        "site": site_id,
        "createDate": current_time,
        "lastUpdate": current_time,
        storage_type: [
            {
                "model": storage_model,
                "serial": storage_serial,
                connection_type: [
                    {
                        task_key_name: {
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

    for attempt in range(max_retries):
        try:
            response = requests.post(
                update_url,
                headers=HEADERS,
                auth=get_awsauth(),
                timeout=30,
                data=json.dumps({"script": script, "upsert": upsert_doc}),
            )
            logger.info("OpenSearch payload: %s", response.request.body)

            if response.status_code in (200, 201):
                logger.info(f"Successfully processed site {site_id}")
                return response.json()

            elif response.status_code == 409:
                logger.warning(
                    f"Version conflict on attempt {attempt + 1} for site {site_id}. Retrying..."
                )
                time.sleep(retry_delay)
                retry_delay *= 2
            else:
                raise Exception(f"OpenSearch error ({response.status_code}): {response.text}")

        except Exception as e:
            logger.error(f"Attempt {attempt + 1} failed for site {site_id}: {str(e)}")
            time.sleep(retry_delay)
            retry_delay *= 2

    raise Exception(f"Failed to update site {site_id} after {max_retries} attempts.")

def lambda_handler(event, context):
    sqs_client = get_sqs_client()
    for record in event["Records"]:
        try:
            body = json.loads(record["body"])
            site_id = body.get("site")
            logger.info(f"Processing record: {body}")

            update_or_create_document(site_id, body)

            receipt_handle = record["receiptHandle"]
            sqs_client.delete_message(QueueUrl=SQS_QUEUE_URL, ReceiptHandle=receipt_handle)
            logger.info(f"Deleted message: {receipt_handle}")

        except Exception as e:
            logger.error(f"Error processing message: {str(e)}", exc_info=True)
            continue

    return {"statusCode": 200, "body": json.dumps({"message": "Processed batch"})}
