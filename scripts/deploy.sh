#!/usr/bin/env bash
set -e

function error {
  echo -e "\n\x1B[31m$1\x1b[0m"
  exit 1
}

function info {
  echo -e "\n\x1B[34m$1\x1b[0m"
}

DEV_HOST="ubuntu@dev.borngosugaming.com"
PROD_HOST="ubuntu@borngosugaming.com"
SERVICE_NAME="bg-mentor"

ENV=$1

if [ "${ENV}" != "dev" ] && [ "${ENV}" != "production" ]; then
    echo "First arg must be the environment, one of 'dev' or 'production'"
    exit 1
fi

if [ "${ENV}" == "production" ]; then
    HOST="${PROD_HOST}"
fi

if [ "${ENV}" == "dev" ]; then
    HOST="${DEV_HOST}"
fi

info "deploying to '${ENV}' (${HOST})"

info "cleanup past deploy failures"
ssh -T -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<-EOF
rm -rf ~/artifact
rm -rf ~/files
EOF

info "copy over new artifacts and files to server"
scp -i .secrets-decrypted/${ENV}/deploy-key.pem -r artifact ${HOST}:~/artifact
scp -i .secrets-decrypted/${ENV}/deploy-key.pem -r deploy/files ${HOST}:~/files

info "stop service, and remove old artifact"
ssh -T -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<-EOF
if [[ "\$(systemctl status | grep ${SERVICE_NAME}.service | grep -v grep)" != "" ]]; then
    sudo systemctl stop ${SERVICE_NAME}
fi
sudo rm -rf /opt/${SERVICE_NAME}
EOF

info "copy over new artifact, and any file updates"
ssh -T -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<-EOF
sudo mkdir -p /opt/${SERVICE_NAME}/
sudo cp -r ~/artifact/${ENV}/* /opt/${SERVICE_NAME}/
sudo cp -r ~/files/all/* /
sudo cp -r ~/files/${ENV}/* /
EOF

info "cleanup after deploy"
ssh -T -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<-EOF
rm -rf ~/artifact
rm -rf ~/files
EOF

info "tell systemd to keep our service up if it doesn't already"
ssh -T -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<-EOF
if [[ "\$(systemctl status | grep ${SERVICE_NAME}.service | grep -v grep)" == "" ]]; then
    sudo systemctl enable ${SERVICE_NAME}
fi
EOF

info "start our service"
ssh -T -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<-EOF
sudo systemctl start ${SERVICE_NAME}
EOF
