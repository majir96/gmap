# Stage 1: Build the executable
FROM golang:1.20-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gmap main.go

# Stage 2: Build a small image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/gmap .

# Command to run the executable
CMD ["./gmap"]