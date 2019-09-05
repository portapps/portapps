#!/bin/bash
set -e

export TOOLS_PATH=${TOOLS_PATH:-/c/portapps_opt}

export ANT_HOME=${ANT_HOME:-$TOOLS_PATH/ant}
export ANT_VERSION=${ANT_VERSION:-1.10.7}

cd ${TOOLS_PATH}

if [ ! -f "${TOOLS_PATH}/apache-ant-${ANT_VERSION}-bin.zip" ]; then
  echo "Downloading http://mirrors.standaloneinstaller.com/apache//ant/binaries/apache-ant-${ANT_VERSION}-bin.zip..."
  curl -fsS -o apache-ant-${ANT_VERSION}-bin.zip http://mirrors.standaloneinstaller.com/apache//ant/binaries/apache-ant-${ANT_VERSION}-bin.zip
  rm -rf ${ANT_HOME}
fi
if [ ! -d "${ANT_HOME}" ]; then
  echo "Extracting apache-ant-${ANT_VERSION}-bin.zip..."
  7z x apache-ant-${ANT_VERSION}-bin.zip -y -o${TOOLS_PATH}/
  mv ${TOOLS_PATH}/apache-ant-${ANT_VERSION} ${ANT_HOME}
fi

ant -version
