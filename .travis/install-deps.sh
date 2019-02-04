#!/bin/bash
set -e

echo "#### Init"
rm -rf /c/portapps_opt
mkdir -p /c/portapps_opt
export PATH=${JAVA_HOME}/bin:${GOPATH}/bin:${GOPATH}/bin/windows_386:${GOROOT}/bin:${NODEJS_HOME}:${PATH}

echo "#### Installing Java JDK and Apache ANT"
choco install jdk8 -params 'installdir=c:\\portapps_opt\\jdk' -y
choco install ant --ignore-dependencies -y

echo "#### Installing Node.js ${NODEJS_VERSION}"
curl -fsS -o node-v${NODEJS_VERSION}-win-x64.7z https://nodejs.org/dist/v${NODEJS_VERSION}/node-v${NODEJS_VERSION}-win-x64.7z
7z x node-v${NODEJS_VERSION}-win-x64.7z -y -o/c/portapps_opt/
mv /c/portapps_opt/node-v${NODEJS_VERSION}-win-x64 ${NODEJS_HOME}
npm config set prefix "c:\\portapps_opt\\nodejs"

echo "#### Installing Go ${GOVERSION}"
curl -fsS -o go${GOVERSION}.windows-amd64.zip https://storage.googleapis.com/golang/go${GOVERSION}.windows-amd64.zip
7z x go${GOVERSION}.windows-amd64.zip -y -o/c/portapps_opt/

echo "#### Finalize"
java -version
ant -version
go version
go env
node --version
npm --version
