# Use the official Go image as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container's working directory
COPY vault_wrapper/ .

# Build the Go program
RUN go build -o /aws

# Code file to execute when the Docker container starts up
ENTRYPOINT ["/aws"]
