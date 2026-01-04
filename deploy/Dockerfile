FROM golang:1.22-alpine AS builder

ARG SERVICE_NAME
ARG SERVICE_PATH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/${SERVICE_NAME} ${SERVICE_PATH}

FROM alpine:latest
ARG SERVICE_NAME

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/${SERVICE_NAME} /app/server

CMD ["./server", "-f", "etc/config.yaml"]
