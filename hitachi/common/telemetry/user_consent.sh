#!/bin/bash

ROOT_DIR="/opt/hitachi/terraform"
CONFIG_FILE="$ROOT_DIR/config.json"
CONSENT_FILE="$ROOT_DIR/user_consent.json"
TELEMETRY_DIR="$ROOT_DIR/telemetry"

# Ensure telemetry directory exists
mkdir -p "$TELEMETRY_DIR"

# Read config.json
if [ ! -f "$CONFIG_FILE" ]; then
  echo "Error: config.json not found at $CONFIG_FILE"
  exit 1
fi

CONSENT_MESSAGE=$(jq -r '.user_consent_message' "$CONFIG_FILE")

echo ""
echo "==================== USER CONSENT ===================="
echo "$CONSENT_MESSAGE"
echo "======================================================"
echo ""

# Prompt for user input
read -p "Do you consent to the collection of usage data? (Yes/No): " USER_INPUT
USER_INPUT_LOWER=$(echo "$USER_INPUT" | tr '[:upper:]' '[:lower:]')

if [[ "$USER_INPUT_LOWER" != "yes" && "$USER_INPUT_LOWER" != "no" ]]; then
  echo "Invalid input. Please enter 'Yes' or 'No'."
  exit 1
fi

USER_CONSENT_BOOL=false
if [ "$USER_INPUT_LOWER" == "yes" ]; then
  USER_CONSENT_BOOL=true
fi

# Get current UTC timestamp
CURRENT_TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Initialize or update consent file
if [ -f "$CONSENT_FILE" ]; then
  # Read existing data
  SITE_ID=$(jq -r '.site_id' "$CONSENT_FILE")
  PREVIOUS_CONSENT=$(jq -r '.user_consent_accepted' "$CONSENT_FILE")
  PREVIOUS_TIME=$(jq -r '.time' "$CONSENT_FILE")
  CONSENT_HISTORY=$(jq -c '.consent_history' "$CONSENT_FILE")

  # Append previous record to history
  UPDATED_HISTORY=$(echo "$CONSENT_HISTORY" | jq \
    --argjson prev "{\"user_consent_accepted\": $PREVIOUS_CONSENT, \"time\": \"$PREVIOUS_TIME\"}" \
    '. + [$prev]')
else
  SITE_ID=$(uuidgen)
  UPDATED_HISTORY="[]"
fi

# Write updated consent file
jq -n \
  --arg site_id "$SITE_ID" \
  --argjson user_consent_accepted "$USER_CONSENT_BOOL" \
  --arg time "$CURRENT_TIMESTAMP" \
  --argjson consent_history "$UPDATED_HISTORY" \
  '{
    site_id: $site_id,
    user_consent_accepted: $user_consent_accepted,
    time: $time,
    consent_history: $consent_history
  }' > "$CONSENT_FILE"

echo ""
echo "✅ User consent has been recorded successfully."
