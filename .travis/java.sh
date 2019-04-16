#!/bin/bash
set -e

#choco install jdk8 -params 'installdir=c:\\portapps_opt\\jdk' -y
curl -fsS -o jdk-${JAVA_VERSION}-windows-i586.zip https://dl.bintray.com/crazy/tools/jdk-${JAVA_VERSION}-windows-i586.zip
7z x jdk-${JAVA_VERSION}-windows-i586.zip -y -o/c/portapps_opt/

java -version
