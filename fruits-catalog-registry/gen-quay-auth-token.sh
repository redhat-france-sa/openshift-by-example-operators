#!/bin/bash
set -e

QUAY_URL=https://quay.io
#QUAY_URL=https://quay-2884.apps.shared-na4.na4.openshift.opentlc.com

echo -n "Username: "
read USERNAME
echo -n "Password: "
read -s PASSWORD 
echo

AUTH_TOKEN=$(curl -sH "Content-Type: application/json" \
    -XPOST $QUAY_URL/cnr/api/v1/users/login \
    -d '{"user": {"username": "'"${USERNAME}"'", "password": "'"${PASSWORD}"'"}}' | jq -r '.token')

echo $AUTH_TOKEN