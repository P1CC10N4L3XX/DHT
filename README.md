# 🧩 Distributed Hash Table (DHT) System

Questo progetto implementa una **Distributed Hash Table (DHT)** basata sulla teoria degli **alberi binari**, capace di gestire in modo efficiente le operazioni di:

- 🔗 **JOIN** – Inserimento di nuovi nodi nella rete  
- ❌ **LEAVE** – Uscita controllata di un nodo con ridistribuzione delle risorse  
- 💾 **PUT** – Inserimento di una risorsa nella DHT  
- 📥 **GET** – Recupero di una risorsa dalla DHT  

Ogni nodo è eseguito come **container Docker** e comunica con gli altri attraverso una rete bridge, formando una rete distribuita autonoma e scalabile.

---

## ⚙️ Requisiti

Assicurati di avere installato i seguenti strumenti sul tuo sistema:

- 🐋 [Docker](https://docs.docker.com/get-docker/)
- 🧱 [Docker Compose](https://docs.docker.com/compose/install/)
- 🖥️ [tmux](https://github.com/tmux/tmux/wiki) (necessario per avviare i nodi in sessioni parallele)

---

## Build dell’immagine Docker

Per costruire l’immagine del nodo DHT:

```bash
docker build -t dht-node .
```

Questo comando creerà un’immagine locale chiamata **`dht-node`**, che rappresenta un singolo nodo della rete distribuita.

---
##  Setup Docker compose

Per generare il file docker-compose.yml

```bash
./generate_yml.sh
```


---

## Avvio della rete DHT

Per avviare tutti i container (nodi della rete):

```bash
docker compose up -d
```

Questo comando esegue tutti i servizi definiti nel file `docker-compose.yml` in modalità **detached**, avviando una rete distribuita composta da più nodi.

---

## Esecuzione del sistema DHT

Una volta che i container sono attivi, puoi avviare i processi dei nodi con:

```bash
./start.sh
```

> ⚠️ **Nota:**  
> Lo script `start.sh` utilizza **tmux** per avviare sessioni parallele, una per ciascun nodo.  
> Assicurati che `tmux` sia installato prima di eseguire questo comando.

---


##  Pulizia dell’ambiente

Per arrestare e rimuovere tutti i container:

```bash
docker compose down
```

Per rimuovere anche i dati salvati dai nodi:

```bash
rm -rf ./data
```


## Autore

**Alessandro Piccione**  
📚 Progetto sviluppato per il corso di Sistemi Distribuiti e Cloud Computing  
🏫 Università degli studi di Roma "Tor Vergata"


