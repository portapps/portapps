#!/bin/bash
set -e

export TOOLS_PATH=${TOOLS_PATH:-/c/portapps_opt}

export JAVA_HOME=${JAVA_HOME:-$TOOLS_PATH/jdk}
export JAVA_VERSION=${JAVA_VERSION:-11.0.7+10}
export JAVA_VERSION_PATCH=${JAVA_VERSION_PATCH:-.2}
export OPENJDK_ZIP_URL=${OPENJDK_ZIP_URL:-https://github.com/AdoptOpenJDK/openjdk11-binaries/releases/download/jdk-${JAVA_VERSION}${JAVA_VERSION_PATCH}/OpenJDK11U-jdk_x64_windows_hotspot_${JAVA_VERSION//+/_}.zip}
OPENJDK_ZIP_FILENAME=${OPENJDK_ZIP_URL##*/}

cd ${TOOLS_PATH}

if [ ! -f "${OPENJDK_ZIP_FILENAME}" ]; then
  echo "Downloading ${OPENJDK_ZIP_URL}..."
  curl -fsS -o ${OPENJDK_ZIP_FILENAME} "${OPENJDK_ZIP_URL}"
  rm -rf ${JAVA_HOME}
fi
if [ ! -d "${JAVA_HOME}" ] && [ -f "${OPENJDK_ZIP_FILENAME}" ]; then
  echo "Extracting ${OPENJDK_ZIP_FILENAME}..."
  7z x ${OPENJDK_ZIP_FILENAME} -y -o${TOOLS_PATH}/
  mv ${TOOLS_PATH}/jdk-${JAVA_VERSION} ${JAVA_HOME}
fi

java -version
