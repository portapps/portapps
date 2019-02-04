#!/bin/bash
set -e

curl -fsS -o go${GOVERSION}.windows-amd64.zip https://storage.googleapis.com/golang/go${GOVERSION}.windows-amd64.zip
7z x go${GOVERSION}.windows-amd64.zip -y -o/c/portapps_opt/

go version
go env
