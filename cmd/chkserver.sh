#!/bin/bash

contstat="temp"
serverstat="temp"

if [ "$(docker container ls | grep -c Up)" == "4" ]; then
    contstat="$(date) Container Deployment Success";
else
    echo "$(date) Container Deployment Failure";
    docker container ls -a
    exit 1;
fi

if [ "$(curl -s -o /dev/null -I -w "%{http_code}\n" https://api.hology.id)" == "400" ]; then
    serverstat="Server is running OK"
else
    serverstat="Server failed to run"
fi

echo "============================"
echo "STATUS: "
echo "$contstat"
echo "$serverstat"
echo "============================"
