export GOROOT=/opt/go
export PATH="$($GOROOT/bin/go env GOPATH)/bin:$GOROOT/bin:$PATH"

go version
tfplugindocs version

export TF_LOG=DEBUG
export TF_LOG_PATH="terraform_debug.log"
export TF_VAR_hitachi_storage_user=admin
export TF_VAR_hitachi_storage_password=Hitachi1

