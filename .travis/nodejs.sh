#!/bin/bash
set -e

curl -fsS -o node-v${NODEJS_VERSION}-win-x64.7z https://nodejs.org/dist/v${NODEJS_VERSION}/node-v${NODEJS_VERSION}-win-x64.7z
7z x node-v${NODEJS_VERSION}-win-x64.7z -y -o/c/portapps_opt/
mv /c/portapps_opt/node-v${NODEJS_VERSION}-win-x64 ${NODEJS_HOME}
npm config set prefix "c:\\portapps_opt\\nodejs"

node --version
npm --version
