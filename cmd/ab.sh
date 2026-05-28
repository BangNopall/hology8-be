#!/bin/sh

ab -n 10000 -c 1 -p data.json -T application/json \
    -H "x-api-key: Key 987adb66d54716e09086d045ca683f4aea45702067785df61c631ade1d62d9f7" \
    http://localhost:8080/api/v1/users/login