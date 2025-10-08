# Stage di build
FROM golang:1.23-alpine AS build

WORKDIR /app

# Copia moduli e scarica le dipendenze
COPY go.mod go.sum* ./
RUN go mod download || true

# Copia tutto il codice sorgente
COPY . .

# Compila l'applicazione
RUN go build -o /dht .

# Stage finale
FROM alpine:3.20

WORKDIR /app

# Copia l'eseguibile dal build stage
COPY --from=build /dht /app/dht

# Copia eventuali sorgenti se servono runtime (es. routing CSV)
COPY src/routing/ /app/src/routing/
COPY src/Data/ /app/src/Data/

# Apri stdin e tty per poter interagire con l'app via docker exec
# Non avvia l'app direttamente, così eviti conflitti di porta
# L'app verrà avviata manualmente con docker exec
CMD ["/app/dht"]
