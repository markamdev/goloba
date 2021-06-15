# --- build application
FROM golang:1.14 AS builder
WORKDIR /temp

COPY go.mod .
RUN go mod download

COPY . .

RUN go build ./cmd/goloba
RUN go build ./cmd/dummyserver

# --- create image
FROM ubuntu:20.04

WORKDIR /usr/bin
COPY --from=builder /temp/build/goloba .
COPY --from=builder /temp/build/dummyserver .

# Listening port for goloba
EXPOSE 8060
# Listening port for dummyserver
EXPOSE 8070
