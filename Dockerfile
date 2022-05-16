FROM golang:1.17.8-alpine AS builder

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -o zaca .

FROM ubuntu:20.04

WORKDIR /zaca

COPY --from=builder /build/zaca .
COPY --from=builder /build/database/mysql/migrations ./database/mysql/migrations
COPY --from=builder /build/conf.prod.yml .
COPY --from=builder /build/conf.test.yml .
RUN chmod +x capitalizone

# API service
CMD ["./zaca", "api"]

# TLS service
# CMD ["./zaca", "api"]

# OCSP service
# CMD ["./zaca", "api"]
