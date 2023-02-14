FROM golang:alpine as builder
WORKDIR /app
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o app ./cmd/it-news-bot


FROM alpine:latest
WORKDIR /root/
COPY --from=builder ./app/app .
EXPOSE 8080
RUN mkdir data
CMD ["./app"]