#!/bin/bash
set -e

export TOOLS_PATH=${TOOLS_PATH:-/c/portapps_opt}

export GOPATH=${GOPATH:-$TOOLS_PATH/gopath}
export GOROOT=${GOROOT:-$TOOLS_PATH/go}
export GOVERSION=${GOVERSION:-1.12.4}

cd ${TOOLS_PATH}

if [ ! -f "${TOOLS_PATH}/go${GOVERSION}.windows-amd64.zip" ]; then
  echo "Downloading https://storage.googleapis.com/golang/go${GOVERSION}.windows-amd64.zip..."
  curl -fsS -o go${GOVERSION}.windows-amd64.zip https://storage.googleapis.com/golang/go${GOVERSION}.windows-amd64.zip
  rm -rf ${GOPATH} ${GOROOT}
fi
if [ ! -d "${GOROOT}" ]; then
  echo "Extracting go${GOVERSION}.windows-amd64.zip..."
  7z x go${GOVERSION}.windows-amd64.zip -y -o${TOOLS_PATH}/
fi

go version
go env
