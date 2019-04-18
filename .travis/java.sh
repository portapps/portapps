#!/bin/bash
set -e

export TOOLS_PATH=${TOOLS_PATH:-/c/portapps_opt}

export JAVA_HOME=${JAVA_HOME:-$TOOLS_PATH/jdk}
export JAVA_VERSION=${JAVA_VERSION:-11.0.2}

cd ${TOOLS_PATH}

if [ ! -f "${TOOLS_PATH}/jdk-${JAVA_VERSION}-windows-i586.zip" ]; then
  echo "Downloading https://download.java.net/java/GA/jdk11/9/GPL/openjdk-${JAVA_VERSION}_windows-x64_bin.zip..."
  curl -fsS -o openjdk-${JAVA_VERSION}_windows-x64_bin.zip https://download.java.net/java/GA/jdk11/9/GPL/openjdk-${JAVA_VERSION}_windows-x64_bin.zip
  rm -rf ${JAVA_HOME}
fi
if [ ! -d "${JAVA_HOME}" ]; then
  echo "Extracting jdk-${JAVA_VERSION}-windows-i586.zip..."
  7z x openjdk-${JAVA_VERSION}_windows-x64_bin.zip -y -o${TOOLS_PATH}/
  mv ${TOOLS_PATH}/jdk-${JAVA_VERSION} ${JAVA_HOME}
fi

java -version
