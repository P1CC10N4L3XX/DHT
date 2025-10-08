#!/bin/bash

SESSION_NAME="dht_cluster"
TMUX_CONF="/tmp/tmux_dht.conf"

# Legge i servizi definiti nel docker-compose.yml
NODES=$(docker-compose config --services | grep '^node')

if [ -z "$NODES" ]; then
  echo "❌ Nessun nodo trovato nel docker-compose.yml"
  exit 1
fi

# Crea configurazione tmux temporanea
cat > "$TMUX_CONF" <<EOF
set -g mouse on
set -g status on
setw -g window-status-current-style bg=green,fg=black
set -g status-style bg=black,fg=white
EOF

FIRST_NODE=true

for NODE in $NODES; do
  # Controlla se il container è attivo
  RUNNING=$(docker ps -q -f name=^${NODE}$)
  if [ -z "$RUNNING" ]; then
    echo "⚠️  Il container $NODE non è in esecuzione. Avvia con: docker compose up -d $NODE"
    continue
  fi

  # Comando per attach
  CMD="echo Press enter...;docker attach $NODE"

  if [ "$FIRST_NODE" = true ]; then
    tmux -f "$TMUX_CONF" new-session -d -s "$SESSION_NAME" -n "$NODE" "$CMD"
    FIRST_NODE=false
  else
    tmux new-window -t "$SESSION_NAME" -n "$NODE" "$CMD"
  fi
done

# Attacca la sessione tmux
tmux attach -t "$SESSION_NAME"


