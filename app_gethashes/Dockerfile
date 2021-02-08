#
# Build the application
#

# Start from the latest golang base image
FROM golang:latest as builder

# Maintainer info
LABEL maintainer="Hans Åge Martinsen <hamartin@moshwire.com>"

# Setting working directory
WORKDIR /app

# Copy go mod and sum files
COPY src/go.mod src/go.sum src/gethashes.go ./

# Download all dependencies
RUN go mod download

# Build the go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gethashes .

#
# Run the application
#

# We start with the latest alpine base image
FROM alpine:latest

# Update, upgrade and install packages
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash

# Setting working directory
WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/gethashes .

# Copy the static resources
ADD img ./img
ADD templates ./templates

# Expose port 7001 to the outside world
EXPOSE 7001

# Command to run the executable
CMD ["./gethashes"]
