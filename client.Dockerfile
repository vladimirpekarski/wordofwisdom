FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN cd /app/cmd/client && \
    go build -o /client/client


FROM alpine:latest AS production

RUN apk add --no-cache ca-certificates

COPY --from=builder /client /client
COPY --from=builder app/config /config
ENV CONFIG_PATH=config/local.yaml

CMD /client/client

