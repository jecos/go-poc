#!/bin/bash

CONTAINER_NAME="starrocks_test"

# Check if the container is already running
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    echo "Container $CONTAINER_NAME is already running."
else
    # Check if the container exists but is stopped
    if [ "$(docker ps -aq -f status=exited -f name=$CONTAINER_NAME)" ]; then
        echo "Starting existing container $CONTAINER_NAME..."
        docker start $CONTAINER_NAME
    else
        echo "Creating and starting new container $CONTAINER_NAME..."
        docker run -d --name $CONTAINER_NAME -p 9030:9030 -p 8030:8030 -p 8040:8040 starrocks/allin1-ubuntu
    fi

 
fi