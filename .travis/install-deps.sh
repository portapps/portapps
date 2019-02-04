#!/bin/bash
set -e

export GOPATH=${GOPATH:-/c/gopath}
export GOROOT=${GOROOT:-/c/go}
export GOVERSION=${GOVERSION:-1.11.5}
export JAVA_HOME=${JAVA_HOME:-/c/jdk}
export NODEJS_HOME=${NODEJS_HOME:-/c/nodejs}
export PATH=${JAVA_HOME}/bin:${GOPATH}/bin:${GOPATH}/bin/windows_386:${GOROOT}/bin:${NODEJS_HOME}:${PATH}

# Install
rm -rf ${JAVA_HOME} ${GOROOT}
choco install jdk8 -params "installdir=${JAVA_HOME}" -y
choco install ant --ignore-dependencies -y
choco install nodejs.install -ia "INSTALLDIR=${NODEJS_HOME}" -y
curl -fsS -o go${GOVERSION}.windows-amd64.zip https://storage.googleapis.com/golang/go${GOVERSION}.windows-amd64.zip
7z x go${GOVERSION}.windows-amd64.zip -y -o/c/

# Configure / Check
npm config set prefix "${NODEJS_HOME}"
java -version
ant -version
go version
go env
node --version
npm --version
