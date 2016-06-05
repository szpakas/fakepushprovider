#!/usr/bin/env bash

export APP_SERVICE="apns"
export APP_APPS_FILE="../generator/tmp/apns-apps.json"
export APP_INSTANCES_FILE="../generator/tmp/apns-inst.json"
export APP_APNS_CERT_FILE="tmp/self-001-cert.pem"
export APP_APNS_KEY_FILE="tmp/self-001-key.pem"
export APP_LOG_LEVEL="all"

./main
