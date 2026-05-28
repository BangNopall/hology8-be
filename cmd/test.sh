#!/bin/bash

# Usage: 
# 1. cmd/test.sh // running all unit testing in internal directory
# 2. cmd/test.sh user // running all unit testing in internal/user directory, you can change the user entity to other entity that are available in internal dir
# 3. cmd/test.sh user service // running all unit testing in internal/user/service directory, you can also change the entity and the layer

# Summary
# Args 1 : entity
# Args 2 : layer
# Args 3 : Verbose or Not

if [ -n "$1" ]; then 

    if [ -n "$2" ]; then 
        echo "Running all tests in internal/app/$1/$2 directory"
        sleep 1
        res=$(go test ./internal/app/$1/$2/...)
    else 
        echo "Running all tests in internal/app/$1 directory"
        sleep 1
        res=$(go test ./internal/app/$1/...)
    fi 

else 
    echo "Running all tests in internal/app directory"
    sleep 1
    res=$(go test ./internal/app/...) 
fi 

res=$(echo "$res" | grep FAIL)

if [ -z "$res" ]; then 
    echo "All unit tests passed!"
else 
    echo "$res"
    echo "Unit tests failed!"
fi

