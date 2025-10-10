#!/bin/bash

# Controllo argomento
if [ $# -lt 1 ]; then
    echo "Uso: $0 <num_node>"
    exit 1
fi

NUM_NODES=$1
FILE_YML="docker-compose.yml"

# Inizio file YAML
cat > $FILE_YML <<EOL
services:
EOL

# --- node0 (entry node) ---
cat >> $FILE_YML <<EOL
  node0:
    image: dht-node
    container_name: node0
    tty: true
    stdin_open: true
    hostname: node0
    command: ["./dht", "-entry"]
    ports:
      - "50052:50051"
    volumes:
      - /app/data
    environment:
      - NODE_NAME=node0
      - NODE_PORT=50051
      - ENTRY_HOST=node0
      - ENTRY_PORT=50051
    restart: "no"

EOL

# --- altri nodi ---
PORT_BASE=50053
for i in $(seq 1 $((NUM_NODES-1))); do
    PORT=$((PORT_BASE + i - 1))
    cat >> $FILE_YML <<EOL
  node$i:
    image: dht-node
    container_name: node$i
    tty: true
    stdin_open: true
    hostname: node$i
    command: ["./dht"]
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

# --- rete ---
cat >> $FILE_YML <<EOL
networks:
  dht_net:
    driver: bridge
EOL

echo "✅ File ${FILE_YML} generato con $NUM_NODES nodi."

