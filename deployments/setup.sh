#!/bin/bash

set -e

# Define directories
GOPATH=$(go env GOPATH)
FABRIC_DIR="./fabric"
FABRIC_SCRIPT="./install-fabric.sh"
IPFS_CONTAINER_NAME="ipfs-kubo"

# Check prerequisites
if ! command -v docker &>/dev/null || ! command -v docker-compose &>/dev/null || ! command -v go &>/dev/null; then
  echo "Docker, Docker Compose, and Go must be installed."
  exit 1
fi

# Step 1: Install binaries and Docker Images
echo "Installing Hyperledger Fabric binaries and Docker images..."
cd $FABRIC_DIR
$FABRIC_SCRIPT d

# Step 2: Set up the test network
cd "fabric-samples/test-network"
echo "Setting up Fabric Test Network..."
./network.sh down
./network.sh up createChannel -ca

# Step 3: Start the chaincode container
./network.sh deployCC -ccn basic -ccp ../blockchain/chaincode/proteomic -ccl go

# Step 4: Set up IPFS with Kubo
if [ ! "$(docker ps -q -f name=$IPFS_CONTAINER_NAME)" ]; then
  echo "Setting up IPFS Kubo..."
  docker run -d --name $IPFS_CONTAINER_NAME \
    -v kubo_staging:/export \
    -v kubo_data:/data/ipfs \
    -p 4001:4001 -p 5001:5001 -p 8080:8080 \
    ipfs/kubo:latest daemon --init
else
  echo "IPFS Kubo container already running."
fi

# Step 5: Run Microservices with Docker Compose

if [ -f "./docker-compose.yml" ]; then
  docker-compose up -d
  echo "Microservices are running."
else
  echo "docker-compose.yml not found in $MICROSERVICES_DIR"
  exit 1
fi

# Summary
echo "Setup completed successfully!"
echo "- Fabric Test Network is running."
echo "- IPFS Kubo is running as container $IPFS_CONTAINER_NAME."
echo "- Microservices are running using Docker Compose."
