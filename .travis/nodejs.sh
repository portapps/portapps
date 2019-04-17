#!/bin/bash
set -e

export TOOLS_PATH=${TOOLS_PATH:-/c/portapps_opt}

export NODEJS_VERSION=${NODEJS_VERSION:-11.14.0}
export NODEJS_HOME=${NODEJS_HOME:-$TOOLS_PATH/nodejs}

cd ${TOOLS_PATH}

if [ ! -f "${TOOLS_PATH}/go${GOVERSION}.windows-amd64.zip" ]; then
  echo "Downloading https://nodejs.org/dist/v${NODEJS_VERSION}/node-v${NODEJS_VERSION}-win-x64.7z..."
  curl -fsS -o node-v${NODEJS_VERSION}-win-x64.7z https://nodejs.org/dist/v${NODEJS_VERSION}/node-v${NODEJS_VERSION}-win-x64.7z
  rm -rf ${NODEJS_HOME}
fi
if [ ! -d "${NODEJS_HOME}" ]; then
  echo "Extracting node-v${NODEJS_VERSION}-win-x64.7z..."
  7z x node-v${NODEJS_VERSION}-win-x64.7z -y -o${TOOLS_PATH}/
  mv ${TOOLS_PATH}/node-v${NODEJS_VERSION}-win-x64 ${NODEJS_HOME}
fi

npm config set prefix "c:\\portapps_opt\\nodejs"
node --version
npm --version
