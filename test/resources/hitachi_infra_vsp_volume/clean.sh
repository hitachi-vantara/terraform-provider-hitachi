!#/bin/sh

echo "Start cleaning terraform files"
rm -rf .terraform
rm -rf .terraform.lock.hcl
rm -rf terraform.tfstate*
echo "Done"
