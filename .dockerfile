FROM golang:1.25.4 AS build
WORKDIR /app

COPY go.mod go.sum ./
COPY go.mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service_good ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/service_good .
COPY --from=build /app/internal/config ./config

EXPOSE 50001
CMD ["/app/service_good"]