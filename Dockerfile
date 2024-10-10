# Use an official Go image
FROM golang:1.17-alpine

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum, and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Set the PORT environment variable
ENV PORT=8080

# Expose the port that the application will use
EXPOSE 8080

# Command to run the application
CMD ["./main"]
