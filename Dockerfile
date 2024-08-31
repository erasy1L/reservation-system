FROM golang:1.23.0-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine AS hoster

WORKDIR /app

COPY --from=builder /build/app ./app

ENTRYPOINT [ "./app" ]