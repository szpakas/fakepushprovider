#!/usr/bin/env bash

export APP_LOG_LEVEL="all"
export APP_APPS_FILE="data/fcm-apps.json"
export APP_INSTANCES_FILE="data/fcm-inst.json"
export APP_GCM_ENDPOINT="http://localhost:8080"

./main
