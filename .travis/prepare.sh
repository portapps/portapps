#!/bin/bash
set -e

export TOOLS_PATH=${TOOLS_PATH:-/c/portapps_opt}

export JAVA_HOME=${JAVA_HOME:-$TOOLS_PATH/jdk}
export ANT_HOME=${ANT_HOME:-$TOOLS_PATH/ant}
export GOPATH=${GOPATH:-$TOOLS_PATH/gopath}
export GOROOT=${GOROOT:-$TOOLS_PATH/go}
export NODEJS_HOME=${NODEJS_HOME:-$TOOLS_PATH/nodejs}

export PATH=${JAVA_HOME}/bin:${ANT_HOME}/bin:${GOPATH}/bin:${GOPATH}/bin/windows_386:${GOROOT}/bin:${NODEJS_HOME}:${PATH}

mkdir -p ${TOOLS_PATH}

echo "TOOLS_PATH=$TOOLS_PATH"
echo "JAVA_HOME=$JAVA_HOME"
echo "ANT_HOME=$ANT_HOME"
echo "GOPATH=$GOPATH"
echo "GOROOT=$GOROOT"
echo "NODEJS_HOME=$NODEJS_HOME"
echo "PATH=$PATH"
