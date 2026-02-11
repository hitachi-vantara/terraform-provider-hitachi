#!/bin/bash
# Terraform Docs generate script.

#Copy the docs dir before genereating the documentation
tfplugindocs generate

# unalias cp
guide_path="docs/guides/"
mkdir -p ${guide_path}
cp ./README.md ${guide_path}
cp ./ReleaseNotes.md ${guide_path}

