FABRIC_DIR="./fabric"
IPFS_CONTAINER_NAME="ipfs-kubo"

echo "Stopping Fabric Test Network"
cd "$FABRIC_DIR/fabric-samples/test-network" || exit
./network.sh down

echo "Killing IPFS container"
docker stop $IPFS_CONTAINER_NAME
