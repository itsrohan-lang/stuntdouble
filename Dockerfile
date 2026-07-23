# Use the official Golang image to create a build artifact.
FROM golang:1.25 AS builder

WORKDIR /app
COPY . .
WORKDIR /app/cli

# Build the CLI statically for maximum portability
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o stuntdouble main.go

# Use a minimal alpine image for the final container
FROM alpine:3.19

RUN apk --no-cache add ca-certificates docker-cli iptables

WORKDIR /root/
COPY --from=builder /app/cli/stuntdouble .

# Set up StuntDouble directories
RUN mkdir -p /root/.stuntdouble/plugins

ENTRYPOINT ["./stuntdouble"]
CMD ["serve"]
