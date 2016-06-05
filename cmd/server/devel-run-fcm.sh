#!/usr/bin/env bash

export APP_SERVICE="fcm"
export APP_APPS_FILE="../generator/tmp/fcm-apps.json"
export APP_INSTANCES_FILE="../generator/tmp/fcm-inst.json"
export APP_LOG_LEVEL="all"

./main
