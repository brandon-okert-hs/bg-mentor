#!/bin/bash

DEV_HOST="ubuntu@bg-dev.brandonokert.com"
PROD_HOST="ubuntu@bg.brandonokert.com"
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

ssh -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<'EOF'
echo ">>>cleanup past deploy failures"
rm -rf ~/artifact
rm -rf ~/files
EOF

echo -e "\n>>>copy over new artifacts and files to server"
scp -i .secrets-decrypted/${ENV}/deploy-key.pem -r artifact ${HOST}:~/artifact
scp -i .secrets-decrypted/${ENV}/deploy-key.pem -r deploy/files ${HOST}:~/files

ssh -i .secrets-decrypted/${ENV}/deploy-key.pem ${HOST} <<'EOF'
echo -e "\n>>>stop service, and remove old artifact"
if [[ "$(systemctl status | grep ${SERVICE_NAME}.service | grep -v grep)" != "" ]]; then
    sudo systemctl stop ${SERVICE_NAME}
fi
sudo systemctl stop ${SERVICE_NAME}
sudo rm -rf /opt/${SERVICE_NAME}

echo -e "\n>>>copy over new artifact, and any file updates"
sudo mkdir -p /opt/${SERVICE_NAME}/
sudo cp -r ~/artifact/* /opt/${SERVICE_NAME}/
sudo cp -r ~/files/* /

echo -e "\n>>>cleanup after deploy"
rm -rf ~/artifact
rm -rf ~/files

echo -e "\n>>>tell systemd to keep our service up if it doesn't already"
if [[ "$(systemctl status | grep ${SERVICE_NAME}.service | grep -v grep)" == "" ]]; then
    sudo systemctl enable ${SERVICE_NAME}
fi

echo -e "\n>>>start our service"
sudo systemctl start ${SERVICE_NAME}
EOF
