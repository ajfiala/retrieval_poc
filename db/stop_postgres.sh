#!/bin/bash

# Set the container name
CONTAINER_NAME="rag-demo-go-postgres"

# Stop the PostgreSQL container
sudo docker stop "$CONTAINER_NAME"

# Remove the PostgreSQL container
sudo docker rm "$CONTAINER_NAME"
