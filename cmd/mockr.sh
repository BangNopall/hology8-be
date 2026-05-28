#!/bin/sh

# Generating mock repository

if [ -z $1 ]; then 
    echo "Please pass service entity"
    echo "Exiting..."
    exit 1
fi 

src="domain/contracts/${1}_contracts.go"
dest="internal/app/${1}/repository/mock/${1}_repository_mock.go"

if [ ! -f "$src" ]; then 
    echo "File $src does not exist"
    echo "Exiting..."
    exit 1
fi 

mockgen -package=repository_mock -source="$src" -destination="$dest"  
