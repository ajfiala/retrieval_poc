#!/bin/bash

# Set the container name
CONTAINER_NAME="rag-demo-go-postgres"

# Set the PostgreSQL version
PG_VERSION="16"

# Set the database name, username, and password
DB_NAME="goragdb"
DB_USER="myuser"
DB_PASSWORD="mypassword"

# Set the network name
# NETWORK_NAME="my-network"

# Start the PostgreSQL container with the pgvector image
sudo docker run -d --name "$CONTAINER_NAME" \
  -e POSTGRES_DB="$DB_NAME" \
  -e POSTGRES_USER="$DB_USER" \
  -e POSTGRES_PASSWORD="$DB_PASSWORD" \
  -p 5432:5432 \
  pgvector/pgvector:pg16
