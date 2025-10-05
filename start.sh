#!/bin/bash

# Numero di nodi da avviare
NUM_NODES=3
BASE_PORT=50051

for i in $(seq 0 $((NUM_NODES-1))); do
  NODE_NAME="node$i"
  HOST_PORT=$((BASE_PORT + i))

  # Avvia container se non esiste
  if [ ! "$(docker ps -aq -f name=$NODE_NAME)" ]; then
    docker run --rm -dit --name $NODE_NAME \
      --network dht_net \
      -v /app/data \
      -e NODE_NAME=$NODE_NAME \
      -e NODE_PORT=$HOST_PORT \
      $( [ $i -eq 0 ] && echo "-e ENTRY_HOST=" ) \
      dht-node tail -f /dev/null
  fi

  # Apri terminale e avvia l'app
  if [ $i -eq 0 ]; then
    # Node0 con -entry
    osascript -e "tell application \"Terminal\" to do script \"docker exec -it $NODE_NAME /app/dht -entry\""
  else
    osascript -e "tell application \"Terminal\" to do script \"docker exec -it $NODE_NAME /app/dht\""
  fi
done
