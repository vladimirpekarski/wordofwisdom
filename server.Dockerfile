FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN cd /app/cmd/server && \
    go build -o /srv/server


FROM alpine:latest AS production

RUN apk add --no-cache ca-certificates

COPY --from=builder /srv /srv
COPY --from=builder app/config /config
ENV CONFIG_PATH=config/local.yaml

CMD /srv/server


