FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/runner
RUN make
