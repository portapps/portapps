#!/bin/bash
set -e

export GOPATH=${GOPATH:-/c/portapps_opt/gopath}
export GOROOT=${GOROOT:-/c/portapps_opt/go}
export JAVA_HOME=${JAVA_HOME:-/c/portapps_opt/jdk}
export NODEJS_HOME=${NODEJS_HOME:-/c/portapps_opt/nodejs}
export GOVERSION=${GOVERSION:-1.11.5}
export NODEJS_VERSION=${NODEJS_VERSION:-11.9.0}
export PATH=${JAVA_HOME}/bin:${GOPATH}/bin:${GOPATH}/bin/windows_386:${GOROOT}/bin:${NODEJS_HOME}:${PATH}

rm -rf /c/portapps_opt
mkdir -p /c/portapps_opt
