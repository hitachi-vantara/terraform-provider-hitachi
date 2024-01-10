#!/bin/bash
# Terraform Docs generate script.

#Copy the docs dir before genereating the documentation
guide_path_tmp="/tmp/docs/guides/"
guide_path="docs/guides/"
mkdir -p ${guide_path_tmp}

unalias cp

cp -rf ${guide_path}/* ${guide_path_tmp}

tfplugindocs generate


mkdir -p ${guide_path}

cp -rf ${guide_path_tmp}/* ${guide_path}
