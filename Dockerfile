FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

RUN go build -o /dht .

FROM alpine:3.20
WORKDIR /app

COPY --from=build /dht /app/dht
COPY src/routing/ /app/src/routing/


CMD ["/app/dht"]
