#!/bin/bash
set -e

export TOOLS_PATH=${TOOLS_PATH:-/c/portapps_opt}

export JAVA_HOME=${JAVA_HOME:-$TOOLS_PATH/jdk}
export JAVA_VERSION=${JAVA_VERSION:-8u211}

cd ${TOOLS_PATH}

if [ ! -f "${TOOLS_PATH}/jdk-${JAVA_VERSION}-windows-i586.zip" ]; then
  echo "Downloading https://dl.bintray.com/crazy/tools/jdk-${JAVA_VERSION}-windows-i586.zip..."
  wget -q https://dl.bintray.com/crazy/tools/jdk-${JAVA_VERSION}-windows-i586.zip
  rm -rf ${JAVA_HOME}
fi
if [ ! -d "${JAVA_HOME}" ]; then
  echo "Extracting jdk-${JAVA_VERSION}-windows-i586.zip..."
  7z x jdk-${JAVA_VERSION}-windows-i586.zip -y -o${TOOLS_PATH}/
fi

java -version
