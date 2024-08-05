#!/bin/bash

init_ipfs_node() {
  IPFS_PATH=$1
  API_PORT=$2
  GATEWAY_PORT=$3
  SWARM_TCP_PORT=$4
  SWARM_UDP_PORT=$5

  export IPFS_PATH=$IPFS_PATH

  # Initialize the IPFS node
  ipfs init

  # Configure the ports
  ipfs config Addresses.API /ip4/127.0.0.1/tcp/$API_PORT
  ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/$GATEWAY_PORT
  ipfs config Addresses.Swarm "[\"/ip4/0.0.0.0/tcp/$SWARM_TCP_PORT\", \"/ip4/0.0.0.0/udp/$SWARM_UDP_PORT/quic\"]"
  
  # Start the IPFS node
  ipfs daemon &
}

# Initialize and start the first IPFS node
# init_ipfs_node ~/.ipfs1 5001 8080 4001 4001

# Initialize and start the second IPFS node
# init_ipfs_node ~/.ipfs2 5002 8081 4002 4002

# Initialize and start the third IPFS node
init_ipfs_node ~/.ipfs3 5003 8082 4003 4003