#!/bin/bash

JARFILE=openapi-generator-cli.jar
DOWNLOAD_LOCATION=https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/5.4.0/openapi-generator-cli-5.4.0.jar

if test -f "$JARFILE"; then
    echo "$JARFILE exists."
else
    echo "$JARFILE does not exist, download it from $DOWNLOAD_LOCATION"
    wget $DOWNLOAD_LOCATION -O "$JARFILE"
fi

java -jar $JARFILE generate -i openapi.yaml \
    -g go-server \
    -c generator-config.yaml \
    -o pkg/server \
    -t openapi-templates
