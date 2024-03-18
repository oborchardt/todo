FROM golang:1.22-alpine3.18 AS builder

# gcc required for go-sqlite3 package
RUN apk add build-base
ENV CGO_ENABLED=1

WORKDIR /go/src/todo

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/todo main.go

FROM alpine:3.18

RUN apk add sqlite

WORKDIR /app

# set up database
COPY todo.sql .
RUN sqlite3 todo.db ".read todo.sql"

# copy go build
COPY --from=builder /go/bin/todo ./

EXPOSE 8080

ENTRYPOINT ["./todo"]