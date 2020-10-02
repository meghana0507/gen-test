# Start Golang base image
FROM golang:alpine as builder

# Set Go Environment
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Add Maintainer info
LABEL maintainer="Meghana Kothaluri <kothaluri07@gmail.com>"

# Install git for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download 

# Copy the source from current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN go build -o main .

# Start new stage from the scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Set the current working directory inside the container 
WORKDIR /root/

# Copy pre-built binary files from previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .       

# Expose port 5000 to the outside world
EXPOSE 5000

# Command to run the executable
CMD ["./main"]