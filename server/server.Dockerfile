FROM golang:1.24 AS builder

WORKDIR /app

COPY ./server/ .
RUN go mod download
RUN go mod tidy

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM scratch
COPY --from=builder /app/server /server
COPY --from=builder /app/.env /.env

CMD ["/server"]
