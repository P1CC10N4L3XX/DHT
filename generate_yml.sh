#!/bin/bash

# Controllo argomento
if [ $# -lt 1 ]; then
    echo "Uso: $0 <num_node>"
    exit 1
fi

NUM_NODES=$1
FILE_YML="docker-compose.yml"

# Inizio file yml
cat > $FILE_YML <<EOL
services:
EOL

# Genera node0 (entry node)
cat >> $FILE_YML <<EOL
  node0:
    image: dht-node
    container_name: node0
    hostname: node0
    stdin_open: true
    tty: true
    command: tail -f /dev/null
    ports:
      - "50052:50051"
    volumes:
      - /app/data
    environment:
      - NODE_NAME=node0
      - NODE_PORT=50051
    restart: "no"

EOL

# Genera gli altri nodi
PORT_BASE=50053
for i in $(seq 1 $((NUM_NODES-1))); do
    PORT=$((PORT_BASE + i - 1))
    cat >> $FILE_YML <<EOL
  node$i:
    image: dht-node
    container_name: node$i
    hostname: node$i
    stdin_open: true
    tty: true
    command: tail -f /dev/null
    ports:
      - "${PORT}:50051"
    volumes:
      - /app/data
    environment:
      - ENTRY_HOST=node0
      - ENTRY_PORT=50051
    depends_on:
      - node0
    restart: "no"

EOL
done

# Rete
cat >> $FILE_YML <<EOL
networks:
  dht_net:
    driver: bridge
EOL

echo "docker-compose.yml generato con $NUM_NODES nodi."
