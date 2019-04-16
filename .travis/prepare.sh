#!/bin/bash
set -e

export JAVA_VERSION=${JAVA_VERSION:-8u211}
export GOVERSION=${GOVERSION:-1.12.4}
export NODEJS_VERSION=${NODEJS_VERSION:-11.14.0}

export GOPATH=${GOPATH:-/c/portapps_opt/gopath}
export GOROOT=${GOROOT:-/c/portapps_opt/go}
export JAVA_HOME=${JAVA_HOME:-/c/portapps_opt/jdk}
export NODEJS_HOME=${NODEJS_HOME:-/c/portapps_opt/nodejs}

export PATH=${JAVA_HOME}/bin:${GOPATH}/bin:${GOPATH}/bin/windows_386:${GOROOT}/bin:${NODEJS_HOME}:${PATH}

rm -rf /c/portapps_opt
mkdir -p /c/portapps_opt
