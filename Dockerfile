# Start from the latest golang base image
FROM golang:latest AS build

# Add Maintainer Info
LABEL maintainer="xiaotianxt <tianyp@pku.edu.cn>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o main .

# Start from a scratch image
FROM scratch

# Copy the binary from the build stage
COPY --from=build /app/main /main

# Expose port 8088 to the outside world
EXPOSE 8088

# Command to run the executable
CMD ["/main"]
