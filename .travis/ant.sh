#!/bin/bash
set -e

choco install ant --ignore-dependencies -y

ant -version
