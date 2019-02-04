#!/bin/bash
set -e

choco install jdk8 -params 'installdir=c:\\portapps_opt\\jdk' -y

java -version
