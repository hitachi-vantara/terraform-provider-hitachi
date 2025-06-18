#!/bin/bash

# ------------------------------------------------------------------------------
# Hitachi Terraform Log Bundle Collector
#
# This script collects Terraform-related configuration files, logs, and system
# information to help with support and debugging.
#
# You can run this script directly:
#     bash /opt/hitachi/terraform/bin/logbundle.sh
#
# Usage:
#   ./logbundle.sh [tf_dir1 tf_dir2 ...]
#       One or more directories to scan for .tf/.tfvars files.
#       If no directories are provided as arguments, you will be prompted to enter
#       directories interactively. Press Enter without input to use the default directories:
#       default dirs = "." "/opt/hitachi/terraform/examples"
#
# Environment Variable:
#   TF_MAX_LOGBUNDLES - Maximum number of log bundles to keep (default: 3)
#
# Note:
#   If you're running this inside a Terraform directory and do not specify any
#   directories, only the default paths will be scanned.
#   Be sure to include all relevant directories explicitly if needed.
# ------------------------------------------------------------------------------

set -euo pipefail

# === Start logging to a temp log file and tee to screen ===
SCRIPT_LOG="/tmp/logbundle_output_$(date +%Y%m%d_%H%M%S).log"
exec > >(tee -a "$SCRIPT_LOG") 2>&1

# === Config ===
TF_MAX_LOGBUNDLES="${TF_MAX_LOGBUNDLES:-3}"
echo "Max log bundles: $TF_MAX_LOGBUNDLES (can be set via environment variable TF_MAX_LOGBUNDLES)"

DEFAULT_TF_DIRS=("." "/opt/hitachi/terraform/examples")

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BUNDLE_DIR="/tmp/hitachi_terraform_logbundle_$TIMESTAMP"
ARCHIVE_NAME="hitachi_terraform_logbundle-$TIMESTAMP.tar.gz"
ARCHIVE_OUTPUT_DIR="/opt/hitachi/terraform/logbundles"
FINAL_ARCHIVE="$ARCHIVE_OUTPUT_DIR/$ARCHIVE_NAME"

PLUGIN_PATH="$HOME/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.1/linux_amd64/terraform-provider-hitachi"
INTERNAL_CONFIG="/opt/hitachi/terraform/bin/.internal_config"
USER_CONSENT="/opt/hitachi/terraform/user_consent.json"
TELEMETRY_DIR="/opt/hitachi/terraform/telemetry"
LOG_DIR="/var/log/hitachi/terraform"
INSTALL_LOG="$LOG_DIR/hitachi_terraform_install.log"
UNINSTALL_LOG="$LOG_DIR/hitachi_terraform_uninstall.log"
MAIN_LOG="$LOG_DIR/hitachi-terraform.log"

# === Functions ===

usage() {
    echo "Usage: $0 [tf_dir1 tf_dir2 ...]"
    echo "  Bundle Terraform directories specified as arguments."
    echo "  If no directories are provided, you will be prompted to enter them."
    echo "  Press Enter without input to use default directories: ${DEFAULT_TF_DIRS[*]}"
    echo
    echo "Options:"
    echo "  -h, --help      Show this help message and exit"
    echo
    echo "Environment Variables:"
    echo "  TF_MAX_LOGBUNDLES  Maximum number of log bundles to keep (default: 3)"
    exit 1
}

parse_args() {
    # Check for -h/--help first
    if [[ $# -ge 1 && ("$1" == "-h" || "$1" == "--help") ]]; then
        usage
    fi

    # Disallow any unsupported options (starting with '-')
    for arg in "$@"; do
        if [[ "$arg" == -* ]]; then
            echo "ERROR: Unsupported option '$arg'"
            usage
        fi
    done

    if [ $# -eq 0 ]; then
        echo "No Terraform directories specified as arguments."
        read -rp "Please enter directories to bundle (space separated), or press Enter to use defaults: " user_input

        if [[ -z "$user_input" ]]; then
            echo "Using default directories: ${DEFAULT_TF_DIRS[*]}"
            TF_DIRS=("${DEFAULT_TF_DIRS[@]}")
        else
            # Split user input into array by whitespace
            read -r -a TF_DIRS <<<"$user_input"
        fi
    else
        TF_DIRS=("$@")
    fi
}

check_installed_terraform() {
    if ! command -v terraform &>/dev/null; then
        echo "❌ Terraform is not installed or not in PATH"
        exit 1
    fi
}

install_jq() {
    if ! command -v jq &>/dev/null; then
        echo "❗ jq not found. Attempting to install jq..."
        if [ -f /etc/os-release ] && grep -q "Oracle Linux Server 8" /etc/os-release; then
            if sudo dnf install -y jq &>/dev/null; then
                echo "✅ jq installed successfully."
            else
                echo "❌ Failed to install jq."
                exit 1
            fi
        else
            echo "⚠️ Unsupported OS for auto jq install. Please install jq manually and rerun."
            exit 1
        fi
    fi
}

redact_tf_variables() {
    local input_file="$1"
    local output_file="$2"
    sed -E '
    /variable[[:space:]]+"[^"]+"/,/}/{ 
      s/(default[[:space:]]*=[[:space:]]*)"[^"]*"/\1"<REDACTED>"/g
    }
  ' "$input_file" >"$output_file"
}

collect_machine_info() {
    echo "📦 Collecting machine info..."
    uname -a >"$BUNDLE_DIR/machine_uname.txt"
    env >"$BUNDLE_DIR/env.txt"
    df -h >"$BUNDLE_DIR/disk_usage.txt"
    free -h >"$BUNDLE_DIR/memory.txt" 2>/dev/null || true
    top -b -n 1 | head -n 30 >"$BUNDLE_DIR/top_snapshot.txt" 2>/dev/null || true
}

collect_terraform_info_global() {
    echo "📦 Collecting Terraform version info..."
    terraform version >"$BUNDLE_DIR/terraform_version.txt" 2>&1 || echo "terraform version failed" >"$BUNDLE_DIR/terraform_version.txt"
}

collect_plugin_info() {
    echo "📦 Collecting hitachi terraform plugin version..."

    rpm -qa --qf 'Hitachi Terraform Provider RPM: %{NAME}-%{VERSION}-%{RELEASE}.%{ARCH}\n' | grep HV_Storage_Terraform > "$BUNDLE_DIR/hitachi_terraform_plugin_version.txt"

    if [[ -x "$PLUGIN_PATH" ]]; then
        "$PLUGIN_PATH" -v >>"$BUNDLE_DIR/hitachi_terraform_plugin_version.txt" 2>&1 || echo "Failed to get plugin version" >>"$BUNDLE_DIR/hitachi_terraform_plugin_version.txt"
    else
        echo "Plugin binary not found at $PLUGIN_PATH" >>"$BUNDLE_DIR/hitachi_terraform_plugin_version.txt"
    fi
}

copy_logs_and_configs() {
    echo "📦 Copying hitachi terraform plugin logs..."
    cp -f "$MAIN_LOG" "$BUNDLE_DIR/" 2>/dev/null || true
    cp -f "$INSTALL_LOG" "$BUNDLE_DIR/" 2>/dev/null || true
    cp -f "$UNINSTALL_LOG" "$BUNDLE_DIR/" 2>/dev/null || true

    echo "📦 Copying hitachi terraform plugin config and telemetry files..."
    cp -f "$INTERNAL_CONFIG" "$BUNDLE_DIR/" 2>/dev/null || true
    cp -f "$USER_CONSENT" "$BUNDLE_DIR/" 2>/dev/null || true

    # Copy only telemetry.json and usages.json from telemetry dir
    mkdir -p "$BUNDLE_DIR/telemetry"
    cp -f "$TELEMETRY_DIR/telemetry.json" "$BUNDLE_DIR/telemetry/" 2>/dev/null || true
    cp -f "$TELEMETRY_DIR/usages.json" "$BUNDLE_DIR/telemetry/" 2>/dev/null || true
}

collect_terraform_tf_info() {
    local dir="$1"
    local tf_out_root="$2"

    mkdir -p "$tf_out_root"
    local tf_validate_out="$tf_out_root/terraform_validate.txt"
    local tf_providers_out="$tf_out_root/terraform_providers.txt"
    local tf_show_out="$tf_out_root/terraform_show.json"

    # Always disable color output
    export TF_IN_AUTOMATION=1

    # echo "🔍 Running terraform validate in $dir"
    if ! terraform -chdir="$dir" validate -no-color >"$tf_validate_out" 2>&1; then
        echo "⚠️ validate failed in $dir" >>"$tf_validate_out"
    fi

    # echo "🔍 Running terraform providers in $dir"
    if ! terraform -chdir="$dir" providers -no-color >"$tf_providers_out" 2>&1; then
        echo "⚠️ providers failed in $dir" >>"$tf_providers_out"
    fi

    # echo "🔍 Running terraform show (or plan fallback) in $dir"
    if terraform -chdir="$dir" show -json >"$tf_show_out" 2>/dev/null; then
        jq . "$tf_show_out" >"${tf_show_out}.tmp" && mv "${tf_show_out}.tmp" "$tf_show_out"
    elif terraform -chdir="$dir" plan -out=tfplan >/dev/null 2>&1 &&
        terraform -chdir="$dir" show -json tfplan >"$tf_show_out" 2>/dev/null; then
        jq . "$tf_show_out" >"${tf_show_out}.tmp" && mv "${tf_show_out}.tmp" "$tf_show_out"
        rm -f "$dir/tfplan"
    else
        echo "terraform show/plan failed in $dir" >"$tf_show_out"
    fi
}

has_tf_files() {
    local dir="$1"
    find -L "$dir" -maxdepth 1 -type f \( \
        -name "*.tf" -o \
        -name "*.tfvars" -o \
        -name "*.tfvars.json" -o \
        -name "*.auto.tfvars.json" \
        \) -print -quit | grep -q .
}

collect_valid_tf_dirs() {
    local base_dir="$1"
    local -n valid_dirs_ref=$2 # pass array name by reference

    while IFS= read -r -d '' dir; do
        # Skip hidden dirs except if it is the base_dir itself
        if [[ "$(basename "$dir")" =~ ^\. ]] && [[ "$dir" != "$base_dir" ]]; then
            continue
        fi

        if has_tf_files "$dir"; then
            valid_dirs_ref+=("$dir")
        fi
    done < <(find "$base_dir" -type d -print0)
}

collect_tf_dirs_and_subdirs() {
    mkdir -p "$BUNDLE_DIR/tf_dirs"

    # Normalize all TF_DIRS entries to absolute paths first
    for i in "${!TF_DIRS[@]}"; do
        if [[ -d "${TF_DIRS[i]}" ]]; then
            TF_DIRS[i]="$(cd "${TF_DIRS[i]}" && pwd)"
        else
            echo "⚠️  Warning: ${TF_DIRS[i]} is not a valid directory"
            unset 'TF_DIRS[i]'
        fi
    done

    for base_dir in "${TF_DIRS[@]}"; do
        base_basename="$(basename "$base_dir")"

        # Skip hidden root dir (if needed)
        if [[ "$base_basename" =~ ^\. ]]; then
            echo "Skipping hidden root directory: $base_dir"
            continue
        fi

        echo "🔍 Searching Terraform directories under: $base_dir"
        short_base_dir="$(flatten_path "$base_dir")"
        tfconfig_root="$BUNDLE_DIR/tf_dirs/$short_base_dir"

        valid_tf_dirs=()
        collect_valid_tf_dirs "$base_dir" valid_tf_dirs

        # Remove duplicates from valid_tf_dirs
        declare -A seen_dirs
        unique_valid_tf_dirs=()
        for d in "${valid_tf_dirs[@]}"; do
            if [[ -z "${seen_dirs[$d]+set}" ]]; then
                seen_dirs[$d]=1
                unique_valid_tf_dirs+=("$d")
            fi
        done

        if [[ ${#unique_valid_tf_dirs[@]} -eq 0 ]]; then
            echo "⚠️  Warning: No .tf or .tfvars files found in $base_dir or its subdirectories"
            continue
        fi

        for dir in "${unique_valid_tf_dirs[@]}"; do
            rel_path="${dir#$base_dir}"
            rel_path="${rel_path#/}" # strip leading slash
            out_dir="$tfconfig_root/$rel_path"
            process_tf_dir "$dir" "$out_dir"
        done
    done
}

process_tf_dir() {
    local dir="$1"
    local tf_out_root="$2"

    # Find .tf and .tfvars files directly under dir (no recursion)
    tf_files=()
    while IFS= read -r -d '' tf_file; do
        tf_files+=("$tf_file")
    done < <(find -L "$dir" -maxdepth 1 -type f \( \
        -name "*.tf" -o \
        -name "*.tfvars" -o \
        -name "*.tfvars.json" -o \
        -name "*.auto.tfvars.json" \
        \) -print0)

    # If no Terraform files found, just return (no warning here; warning is done in collect step)
    if [[ ${#tf_files[@]} -eq 0 ]]; then
        return
    fi

    # Prepare output directory for copying redacted files
    mkdir -p "$tf_out_root"

    # Copy and redact .tf and .tfvars files
    for tf_file in "${tf_files[@]}"; do
        rel_path="${tf_file#$dir/}"
        rel_path="${rel_path#/}" # strip leading slash if any
        out_file="$tf_out_root/$rel_path"
        out_dir="$(dirname "$out_file")"
        mkdir -p "$out_dir"
        redact_tf_variables "$tf_file" "$out_file"
    done

    # Copy terraform.tfstate files if present
    find "$dir" -maxdepth 1 -type f -name "terraform.tfstate*" -exec cp -f {} "$BUNDLE_DIR/" \;

    # Collect additional terraform info (terraform commands etc.)
    collect_terraform_tf_info "$dir" "$tf_out_root"

    # Copy TF_LOG_PATH if set
    if [[ -n "${TF_LOG_PATH:-}" && -f "$TF_LOG_PATH" ]]; then
        cp -f "$TF_LOG_PATH" "$BUNDLE_DIR/tf_log_from_env.log" 2>/dev/null || true
    fi

    # all other files in the directory that are not .tf, .tfvars, .json, or hidden files
    # but include logs, shell scripts, and crash files
    # This is useful for collecting additional logs or scripts that might be relevant
    # to the Terraform execution or environment.
    # echo "📦 Copying additional logs and scripts from $dir to $tf_out_root
    find "$dir" -maxdepth 1 -type f ! -name ".*" ! -name "*.tf*" ! -name "*.json" \
        \( -name "*.log" -o -name "*.sh" -o -name "crash*" \) -exec cp -f {} "$tf_out_root/" \;
}

flatten_path() {
    local path="$1"
    IFS='/' read -r -a parts <<<"$path"

    # Filter out empty parts (due to leading '/')
    local non_empty_parts=()
    for part in "${parts[@]}"; do
        [[ -n "$part" ]] && non_empty_parts+=("$part")
    done

    # Take the last 4 parts or fewer
    local last_parts=()
    local start=$(( ${#non_empty_parts[@]} > 4 ? ${#non_empty_parts[@]} - 4 : 0 ))
    for ((i = start; i < ${#non_empty_parts[@]}; i++)); do
        last_parts+=("${non_empty_parts[i]}")
    done

    # Join with underscore
    (IFS=_; echo "${last_parts[*]}")
}

cleanup_old_logbundles() {
    echo "🧹 Cleaning up old logbundles, keeping only the last $TF_MAX_LOGBUNDLES..."
    find "$ARCHIVE_OUTPUT_DIR" -maxdepth 1 -type f -name 'hitachi_terraform_logbundle-*.tar.gz' |
        sort -r |
        awk "NR>$TF_MAX_LOGBUNDLES" |
        while read -r old_file; do
            echo "🗑️  Removing old bundle: $old_file"
            rm -f "$old_file"
        done
}

archive_logbundle() {
    echo "📦 Copying logbundle script output log..."
    echo "📦 Creating logbundle archive at $FINAL_ARCHIVE..."
    cp -f "$SCRIPT_LOG" "$BUNDLE_DIR/logbundle_script_output.log"
    tar -czf "$FINAL_ARCHIVE" -C /tmp "$(basename "$BUNDLE_DIR")" || {
        echo "❌ Failed to create archive at $FINAL_ARCHIVE"
        exit 1
    }
    echo "✅ Log bundle created: $FINAL_ARCHIVE"

    rm -rf "$BUNDLE_DIR" "$SCRIPT_LOG"
    cleanup_old_logbundles
}

# === Main ===

parse_args "$@"
check_installed_terraform
install_jq
mkdir -p "$BUNDLE_DIR"
mkdir -p "$ARCHIVE_OUTPUT_DIR"

collect_terraform_info_global
collect_tf_dirs_and_subdirs
copy_logs_and_configs
collect_plugin_info
collect_machine_info

archive_logbundle
