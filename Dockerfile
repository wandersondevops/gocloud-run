# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19 AS builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Download dependencies.
RUN go mod download

# Copy the source code.
COPY . .

# Build the application.
RUN go build -v -o server

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
FROM debian:stable-slim

# Copy the binary from the builder stage.
COPY --from=builder /app/server /server

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the web service on container startup.
CMD ["/server"]
