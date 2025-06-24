#!/bin/bash

ROOT_DIR="/opt/hitachi/terraform"
CONFIG_FILE="$ROOT_DIR/bin/.internal_config"
CONSENT_FILE="$ROOT_DIR/user_consent.json"
TELEMETRY_DIR="$ROOT_DIR/telemetry"

# Ensure telemetry directory exists
mkdir -p "$TELEMETRY_DIR"

install_jq() {
  if ! command -v jq &>/dev/null; then
    echo "❗ jq not found. Attempting to install jq..."

    if [ -f /etc/os-release ]; then
      if grep -Eq "Oracle Linux|Red Hat Enterprise Linux|CentOS" /etc/os-release; then
        # Use dnf for version 8+, otherwise yum
        if grep -q "VERSION_ID=\"8" /etc/os-release || grep -q "VERSION_ID=8" /etc/os-release; then
          PKG_CMD="dnf"
        else
          PKG_CMD="yum"
        fi

        if sudo "$PKG_CMD" install -y jq &>/dev/null; then
          echo "✅ jq installed successfully."
        else
          echo "❌ Failed to install jq using $PKG_CMD."
          exit 1
        fi
      else
        echo "⚠️ Unsupported OS for automatic jq installation. Please install jq manually and rerun."
        exit 1
      fi
    fi
  fi
}

install_uuidgen() {
  if ! command -v uuidgen &>/dev/null; then
    echo "❗ uuidgen not found. Attempting to install util-linux..."

    if [ -f /etc/os-release ]; then
      if grep -Eq "Oracle Linux|Red Hat Enterprise Linux|CentOS" /etc/os-release; then
        # Use dnf for version 8+, otherwise yum
        if grep -q "VERSION_ID=\"8" /etc/os-release || grep -q "VERSION_ID=8" /etc/os-release; then
          PKG_CMD="dnf"
        else
          PKG_CMD="yum"
        fi

        if sudo "$PKG_CMD" install -y util-linux &>/dev/null; then
          echo "✅ uuidgen installed successfully."
        else
          echo "❌ Failed to install util-linux using $PKG_CMD."
          exit 1
        fi
      else
        echo "⚠️ Unsupported OS for automatic util-linux installation. Please install it manually and rerun."
        exit 1
      fi
    fi
  fi
}

install_jq
install_uuidgen

# Read config
if [ ! -f "$CONFIG_FILE" ]; then
  echo "Error: $CONFIG_FILE not found"
  exit 1
fi

CONSENT_MESSAGE=$(jq -r '.user_consent_message' "$CONFIG_FILE")

echo ""
echo "==================== USER CONSENT ===================="
echo "$CONSENT_MESSAGE"
echo
echo "======================================================"
echo ""

# Prompt for user input
read -p "Do you consent to the collection of usage data? (Yes/No): " USER_INPUT
USER_INPUT_LOWER=$(echo "$USER_INPUT" | tr '[:upper:]' '[:lower:]')

if [[ "$USER_INPUT_LOWER" != "yes" && "$USER_INPUT_LOWER" != "no" ]]; then
  echo "Invalid input. Please enter 'Yes' or 'No'."
  exit 1
fi

# Convert to raw JSON boolean (unquoted)
if [ "$USER_INPUT_LOWER" == "yes" ]; then
  USER_CONSENT_BOOL=true
else
  USER_CONSENT_BOOL=false
fi

# Get current UTC timestamp
CURRENT_TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Initialize or update consent file
if [ -f "$CONSENT_FILE" ]; then
  SITE_ID=$(jq -r '.site_id' "$CONSENT_FILE")
  if [ "$SITE_ID" == "" ]; then
    SITE_ID=$(uuidgen)
  fi
  PREVIOUS_CONSENT=$(jq -r '.user_consent_accepted' "$CONSENT_FILE")
  PREVIOUS_TIME=$(jq -r '.time' "$CONSENT_FILE")
  CONSENT_HISTORY=$(jq -c '.consent_history' "$CONSENT_FILE")

  # Append previous record to history
  UPDATED_HISTORY=$(echo "$CONSENT_HISTORY" | jq \
    --argjson prev "{\"user_consent_accepted\": $PREVIOUS_CONSENT, \"time\": \"$PREVIOUS_TIME\"}" \
    '. + [$prev]')
else
  SITE_ID=$(uuidgen)
  UPDATED_HISTORY='[]'
fi

# Write updated consent file
jq -n \
  --arg site_id "$SITE_ID" \
  --arg time "$CURRENT_TIMESTAMP" \
  --argjson user_consent_accepted $USER_CONSENT_BOOL \
  --argjson consent_history "$UPDATED_HISTORY" \
  '{
    site_id: $site_id,
    user_consent_accepted: $user_consent_accepted,
    time: $time,
    consent_history: $consent_history
  }' >"$CONSENT_FILE"

echo ""
echo "✅ User consent has been recorded successfully."
