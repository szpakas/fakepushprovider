#!/usr/bin/env bash

export APP_LOG_LEVEL="all"

export APP_APPS_FILE="../generator/tmp/apps.json"
export APP_INSTANCES_FILE="../generator/tmp/instances.json"

export APP_GCM_ENDPOINT="http://localhost:8080"

./main
