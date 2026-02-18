FROM golang:1.24-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-planner .


FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/todo-planner .
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/database/scheduler.db

RUN mkdir -p /database

EXPOSE 7540

CMD ["./todo-planner"]