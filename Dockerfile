ARG TAGS=""

FROM golang:1.20-alpine AS builder
ARG TAGS
WORKDIR /app
COPY . .
RUN go build -trimpath -ldflags "-s -w" -tags "$TAGS" -o ./bin/ninedns .

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/bin/ninedns /usr/local/bin/ninedns

VOLUME /app
ENTRYPOINT ["/usr/local/bin/ninedns"]