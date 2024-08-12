# Use the official Golang image to build the app
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the cmd executable
RUN go build -o /app/cmd ./app/service/cmd

# Build the scheduler executable
RUN go build -o /app/scheduler ./app/service/scheduler

# Build the sender executable
RUN go build -o /app/sender ./app/service/sender

# Start a new stage from scratch
FROM debian:stable-slim

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary files from the previous stage
COPY --from=builder /app/cmd .
COPY --from=builder /app/scheduler .
COPY --from=builder /app/sender .

# Create the /etc/calendar directory
RUN mkdir -p /etc/calendar

# Copy the local config.yaml file into the container
COPY config.yaml /etc/calendar/config.yaml

# Copy migrations directory into the container
COPY migrations ./migrations

# Command to run all executables sequentially
CMD ["sh", "-c", "./cmd & ./scheduler & ./sender & wait"]
