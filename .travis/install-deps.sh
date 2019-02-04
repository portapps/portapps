#!/bin/bash
set -e

export GOPATH=${GOPATH:-/c/portapps_opt/gopath}
export GOROOT=${GOROOT:-/c/portapps_opt/go}
export JAVA_HOME=${JAVA_HOME:-/c/portapps_opt/jdk}
export NODEJS_HOME=${NODEJS_HOME:-/c/portapps_opt/nodejs}

export GOVERSION=${GOVERSION:-1.11.5}
export NODEJS_VERSION=${NODEJS_VERSION:-11.9.0}

# init
rm -rf /c/portapps_opt
mkdir -p /c/portapps_opt
export PATH=${JAVA_HOME}/bin:${GOPATH}/bin:${GOPATH}/bin/windows_386:${GOROOT}/bin:${NODEJS_HOME}:${PATH}

# jdk / ant
choco install jdk8 -params 'installdir=c:\\portapps_opt\\jdk' -y
choco install ant --ignore-dependencies -y

# nodejs
curl -fsS -o node-v${NODEJS_VERSION}-win-x64.7z https://nodejs.org/dist/v${NODEJS_VERSION}/node-v${NODEJS_VERSION}-win-x64.7z
7z x node-v${NODEJS_VERSION}-win-x64.7z -y -o/c/portapps_opt/
mv /c/portapps_opt/node-v${NODEJS_VERSION}-win-x64 ${NODEJS_HOME}
npm config set prefix "c:\\portapps_opt\\nodejs"

# go
curl -fsS -o go${GOVERSION}.windows-amd64.zip https://storage.googleapis.com/golang/go${GOVERSION}.windows-amd64.zip
7z x go${GOVERSION}.windows-amd64.zip -y -o/c/portapps_opt/

# check
java -version
ant -version
go version
go env
node --version
npm --version
