FROM golang:1.24 AS builder

WORKDIR /app

COPY ./client1 .
RUN go mod download
RUN go mod tidy

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -o client1 ./cmd/main.go

FROM scratch
COPY --from=builder /app/client1 /client1

CMD ["/client1"]
