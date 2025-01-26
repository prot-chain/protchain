#!/bin/bash

NETWORK_DIR="./fabric_setup/test-network/"
NETWORK_SCRIPT="./fabric_setup/test-network/network.sh"
FABRIC_SCRIPT="./fabric_setup/install-fabric.sh"

set -e

# Define directories
IPFS_CONTAINER_NAME="ipfs-kubo"

# Check prerequisites
if ! command -v docker &>/dev/null || ! command -v docker-compose &>/dev/null || ! command -v go &>/dev/null; then
  echo "Docker, Docker Compose, and Go must be installed."
  exit 1
fi

# Step 1: Install binaries and Docker Images
echo "Installing Hyperledger Fabric Docker images..."
chmod +x $NETWORK_SCRIPT
chmod +x $FABRIC_SCRIPT
$FABRIC_SCRIPT d

# Step 2: Set up the test network
echo "Setting up Fabric Test Network..."
$NETWORK_SCRIPT down
$NETWORK_SCRIPT up createChannel -ca

# Step 3: Start the chaincode container
cd $NETWORK_DIR
./network.sh deployCC -ccn proteomic -ccp ../../../blockchain/chaincode/proteomic -ccl go

# Step 4: Run Microservices and IPFS with Docker Compose
cd ../..
if [ -f "./docker-compose.yml" ]; then
  docker-compose pull && docker-compose up -d
  echo "Microservices are running."
else
  echo "docker-compose.yml not found in ${PWD}"
  exit 1
fi

# Summary
echo "Setup completed successfully!"
echo "- Fabric Test Network is running."
echo "- IPFS Kubo is running as container $IPFS_CONTAINER_NAME."
echo "- Microservices are running using Docker Compose."
