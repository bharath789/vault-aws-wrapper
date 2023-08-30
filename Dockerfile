# Use the official Go image as the base image
FROM golang:1.17-alpine

# Copy the Go source code into the container's working directory
COPY vault_wrapper/* .

# Build the Go program
RUN go build -o aws.go

# Command to run the application
CMD ["./aws.go"]
