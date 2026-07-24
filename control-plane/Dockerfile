FROM golang:1.21-alpine

WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Disable CGO for static binary so we don't need libc
RUN CGO_ENABLED=1 go build -o control-plane main.go

EXPOSE 4439
CMD ["./control-plane"]
